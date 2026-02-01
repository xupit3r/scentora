package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeIngredientRepository struct {
	db *sqlx.DB
}

func NewRecipeIngredientRepository(db *sqlx.DB) *RecipeIngredientRepository {
	return &RecipeIngredientRepository{db: db}
}

// Create adds an ingredient to a recipe version
func (r *RecipeIngredientRepository) Create(ingredient *models.RecipeIngredient) error {
	query := `
		INSERT INTO recipe_ingredients (
			version_id, accord_id, quantity_ml, percentage, notes, created_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, quantity_drops, created_at
	`
	return r.db.QueryRow(
		query,
		ingredient.VersionID,
		ingredient.AccordID,
		ingredient.QuantityMl,
		ingredient.Percentage,
		ingredient.Notes,
	).Scan(&ingredient.ID, &ingredient.QuantityDrops, &ingredient.CreatedAt)
}

// FindByID retrieves an ingredient by ID
func (r *RecipeIngredientRepository) FindByID(id string) (*models.RecipeIngredient, error) {
	var ingredient models.RecipeIngredient
	query := `
		SELECT id, version_id, accord_id, quantity_ml, quantity_drops,
		       percentage, notes, created_at
		FROM recipe_ingredients
		WHERE id = $1
	`
	err := r.db.Get(&ingredient, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ingredient not found")
	}
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// FindByVersionID retrieves all ingredients for a recipe version
func (r *RecipeIngredientRepository) FindByVersionID(versionID string) ([]*models.RecipeIngredient, error) {
	var ingredients []*models.RecipeIngredient
	query := `
		SELECT id, version_id, accord_id, quantity_ml, quantity_drops,
		       percentage, notes, created_at
		FROM recipe_ingredients
		WHERE version_id = $1
		ORDER BY created_at ASC
	`
	err := r.db.Select(&ingredients, query, versionID)
	return ingredients, err
}

// Update updates an ingredient's quantities and notes
func (r *RecipeIngredientRepository) Update(ingredient *models.RecipeIngredient) error {
	query := `
		UPDATE recipe_ingredients
		SET quantity_ml = $1, percentage = $2, notes = $3
		WHERE id = $4
		RETURNING quantity_drops
	`
	err := r.db.QueryRow(
		query,
		ingredient.QuantityMl,
		ingredient.Percentage,
		ingredient.Notes,
		ingredient.ID,
	).Scan(&ingredient.QuantityDrops)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("ingredient not found")
	}
	return err
}

// Delete removes an ingredient from a recipe version
func (r *RecipeIngredientRepository) Delete(id string) error {
	query := `DELETE FROM recipe_ingredients WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("ingredient not found")
	}
	return nil
}

// ExistsInVersion checks if an accord is already in a version
func (r *RecipeIngredientRepository) ExistsInVersion(versionID, accordID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipe_ingredients WHERE version_id = $1 AND accord_id = $2)`
	err := r.db.Get(&exists, query, versionID, accordID)
	return exists, err
}

// GetTotalVolume returns the sum of all ingredient quantities for a version
func (r *RecipeIngredientRepository) GetTotalVolume(versionID string) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(quantity_ml), 0)
		FROM recipe_ingredients
		WHERE version_id = $1
	`
	err := r.db.Get(&total, query, versionID)
	return total, err
}

// CountByVersionID returns the number of ingredients in a version
func (r *RecipeIngredientRepository) CountByVersionID(versionID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM recipe_ingredients WHERE version_id = $1`
	err := r.db.Get(&count, query, versionID)
	return count, err
}

// FindByAccordID retrieves all versions that use a specific accord
func (r *RecipeIngredientRepository) FindByAccordID(accordID string) ([]string, error) {
	var versionIDs []string
	query := `
		SELECT DISTINCT version_id
		FROM recipe_ingredients
		WHERE accord_id = $1
	`
	err := r.db.Select(&versionIDs, query, accordID)
	return versionIDs, err
}
