package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeVersionRepository struct {
	db *sqlx.DB
}

func NewRecipeVersionRepository(db *sqlx.DB) *RecipeVersionRepository {
	return &RecipeVersionRepository{db: db}
}

// Create creates a new recipe version and automatically sets it as active
func (r *RecipeVersionRepository) Create(version *models.RecipeVersion) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the next version number for this recipe
	var nextVersion int
	err = tx.Get(&nextVersion, `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM recipe_versions
		WHERE recipe_id = $1
	`, version.RecipeID)
	if err != nil {
		return err
	}
	version.VersionNumber = nextVersion

	// Create the new version
	query := `
		INSERT INTO recipe_versions (
			recipe_id, version_number, name, notes, is_active, created_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	err = tx.QueryRow(
		query,
		version.RecipeID,
		version.VersionNumber,
		version.Name,
		version.Notes,
		true, // New versions are always active
	).Scan(&version.ID, &version.CreatedAt)
	if err != nil {
		return err
	}

	// Deactivate all other versions for this recipe
	_, err = tx.Exec(`
		UPDATE recipe_versions
		SET is_active = false
		WHERE recipe_id = $1 AND id != $2
	`, version.RecipeID, version.ID)
	if err != nil {
		return err
	}

	// Update the recipe's active_version_id
	_, err = tx.Exec(`
		UPDATE recipes
		SET active_version_id = $1, updated_at = NOW()
		WHERE id = $2
	`, version.ID, version.RecipeID)
	if err != nil {
		return err
	}

	version.IsActive = true
	return tx.Commit()
}

// FindByID retrieves a version by ID
func (r *RecipeVersionRepository) FindByID(id string) (*models.RecipeVersion, error) {
	var version models.RecipeVersion
	query := `
		SELECT id, recipe_id, version_number, name, notes, is_active, created_at
		FROM recipe_versions
		WHERE id = $1
	`
	err := r.db.Get(&version, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("version not found")
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// FindByRecipeID retrieves all versions for a recipe
func (r *RecipeVersionRepository) FindByRecipeID(recipeID string) ([]*models.RecipeVersion, error) {
	var versions []*models.RecipeVersion
	query := `
		SELECT id, recipe_id, version_number, name, notes, is_active, created_at
		FROM recipe_versions
		WHERE recipe_id = $1
		ORDER BY version_number DESC
	`
	err := r.db.Select(&versions, query, recipeID)
	return versions, err
}

// FindActiveByRecipeID retrieves the active version for a recipe
func (r *RecipeVersionRepository) FindActiveByRecipeID(recipeID string) (*models.RecipeVersion, error) {
	var version models.RecipeVersion
	query := `
		SELECT id, recipe_id, version_number, name, notes, is_active, created_at
		FROM recipe_versions
		WHERE recipe_id = $1 AND is_active = true
	`
	err := r.db.Get(&version, query, recipeID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active version found")
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// SetActive sets a version as active and deactivates all others
func (r *RecipeVersionRepository) SetActive(versionID, recipeID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify version exists and belongs to this recipe
	var exists bool
	err = tx.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM recipe_versions 
			WHERE id = $1 AND recipe_id = $2
		)
	`, versionID, recipeID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("version not found or does not belong to this recipe")
	}

	// Deactivate all versions for this recipe
	_, err = tx.Exec(`
		UPDATE recipe_versions
		SET is_active = false
		WHERE recipe_id = $1
	`, recipeID)
	if err != nil {
		return err
	}

	// Activate the specified version
	_, err = tx.Exec(`
		UPDATE recipe_versions
		SET is_active = true
		WHERE id = $1
	`, versionID)
	if err != nil {
		return err
	}

	// Update the recipe's active_version_id
	_, err = tx.Exec(`
		UPDATE recipes
		SET active_version_id = $1, updated_at = NOW()
		WHERE id = $2
	`, versionID, recipeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CountByRecipeID returns the number of versions for a recipe
func (r *RecipeVersionRepository) CountByRecipeID(recipeID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM recipe_versions WHERE recipe_id = $1`
	err := r.db.Get(&count, query, recipeID)
	return count, err
}

// Delete deletes a version (should check if it's not the only version)
func (r *RecipeVersionRepository) Delete(id, recipeID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if this is the only version
	var count int
	err = tx.Get(&count, `SELECT COUNT(*) FROM recipe_versions WHERE recipe_id = $1`, recipeID)
	if err != nil {
		return err
	}
	if count <= 1 {
		return fmt.Errorf("cannot delete the only version of a recipe")
	}

	// Check if this version is active
	var isActive bool
	err = tx.Get(&isActive, `SELECT is_active FROM recipe_versions WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("version not found")
	}
	if err != nil {
		return err
	}

	// If active, make the previous version active
	if isActive {
		var previousVersionID string
		err = tx.Get(&previousVersionID, `
			SELECT id FROM recipe_versions
			WHERE recipe_id = $1 AND id != $2
			ORDER BY version_number DESC
			LIMIT 1
		`, recipeID, id)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE recipe_versions SET is_active = true WHERE id = $1
		`, previousVersionID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE recipes SET active_version_id = $1, updated_at = NOW() WHERE id = $2
		`, previousVersionID, recipeID)
		if err != nil {
			return err
		}
	}

	// Delete the version (cascade will delete ingredients)
	result, err := tx.Exec(`DELETE FROM recipe_versions WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("version not found")
	}

	return tx.Commit()
}
