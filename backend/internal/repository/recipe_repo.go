package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeRepository struct {
	db *sqlx.DB
}

func NewRecipeRepository(db *sqlx.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

// Create creates a new recipe
func (r *RecipeRepository) Create(recipe *models.Recipe) error {
	query := `
		INSERT INTO recipes (
			user_id, name, description, target_volume_ml, status,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		recipe.UserID,
		recipe.Name,
		recipe.Description,
		recipe.TargetVolumeMl,
		recipe.Status,
	).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.UpdatedAt)
}

// FindByID retrieves a recipe by ID
func (r *RecipeRepository) FindByID(id, userID string) (*models.Recipe, error) {
	var recipe models.Recipe
	query := `
		SELECT id, user_id, name, description, target_volume_ml, status,
		       active_version_id, created_at, updated_at
		FROM recipes
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&recipe, query, id, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recipe not found")
	}
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// FindAll retrieves all recipes for a user
func (r *RecipeRepository) FindAll(userID string) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT id, user_id, name, description, target_volume_ml, status,
		       active_version_id, created_at, updated_at
		FROM recipes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&recipes, query, userID)
	return recipes, err
}

// FindByStatus retrieves recipes filtered by status
func (r *RecipeRepository) FindByStatus(userID, status string) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT id, user_id, name, description, target_volume_ml, status,
		       active_version_id, created_at, updated_at
		FROM recipes
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
	`
	err := r.db.Select(&recipes, query, userID, status)
	return recipes, err
}

// Search searches recipes by name or description
func (r *RecipeRepository) Search(userID, searchQuery string) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT id, user_id, name, description, target_volume_ml, status,
		       active_version_id, created_at, updated_at
		FROM recipes
		WHERE user_id = $1 
		  AND (name ILIKE $2 OR description ILIKE $2)
		ORDER BY created_at DESC
	`
	searchPattern := "%" + searchQuery + "%"
	err := r.db.Select(&recipes, query, userID, searchPattern)
	return recipes, err
}

// Update updates an existing recipe
func (r *RecipeRepository) Update(recipe *models.Recipe) error {
	query := `
		UPDATE recipes
		SET name = $1, description = $2, target_volume_ml = $3, 
		    status = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		query,
		recipe.Name,
		recipe.Description,
		recipe.TargetVolumeMl,
		recipe.Status,
		recipe.ID,
		recipe.UserID,
	).Scan(&recipe.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("recipe not found")
	}
	return err
}

// UpdateActiveVersion sets the active version for a recipe
func (r *RecipeRepository) UpdateActiveVersion(recipeID, versionID, userID string) error {
	query := `
		UPDATE recipes
		SET active_version_id = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`
	result, err := r.db.Exec(query, versionID, recipeID, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("recipe not found")
	}
	return nil
}

// Delete deletes a recipe and all associated data (cascades via DB)
func (r *RecipeRepository) Delete(id, userID string) error {
	query := `DELETE FROM recipes WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("recipe not found")
	}
	return nil
}

// FindByTags retrieves recipes that have all specified tags
func (r *RecipeRepository) FindByTags(userID string, tags []string) ([]*models.Recipe, error) {
	if len(tags) == 0 {
		return r.FindAll(userID)
	}
	
	var recipes []*models.Recipe
	// Use a subquery to find recipes that have ALL specified tags
	query := `
		SELECT DISTINCT r.id, r.user_id, r.name, r.description, 
		       r.target_volume_ml, r.status, r.active_version_id,
		       r.created_at, r.updated_at
		FROM recipes r
		WHERE r.user_id = $1
		  AND (
		    SELECT COUNT(DISTINCT rt.tag)
		    FROM recipe_tags rt
		    WHERE rt.recipe_id = r.id AND rt.tag = ANY($2)
		  ) = $3
		ORDER BY r.created_at DESC
	`
	err := r.db.Select(&recipes, query, userID, tags, len(tags))
	return recipes, err
}

// CountByStatus returns count of recipes by status for a user
func (r *RecipeRepository) CountByStatus(userID string) (map[string]int, error) {
	rows, err := r.db.Query(`
		SELECT status, COUNT(*) as count
		FROM recipes
		WHERE user_id = $1
		GROUP BY status
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

// FindByAccordID retrieves all recipes that use a specific accord
func (r *RecipeRepository) FindByAccordID(accordID, userID string) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT DISTINCT r.id, r.user_id, r.name, r.description,
		       r.target_volume_ml, r.status, r.active_version_id,
		       r.created_at, r.updated_at
		FROM recipes r
		INNER JOIN recipe_versions rv ON rv.recipe_id = r.id
		INNER JOIN recipe_ingredients ri ON ri.version_id = rv.id
		WHERE ri.accord_id = $1 AND r.user_id = $2
		ORDER BY r.created_at DESC
	`
	err := r.db.Select(&recipes, query, accordID, userID)
	return recipes, err
}

// FindByCollection retrieves all recipes in a collection
func (r *RecipeRepository) FindByCollection(collectionID, userID string) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT r.id, r.user_id, r.name, r.description,
		       r.target_volume_ml, r.status, r.active_version_id,
		       r.created_at, r.updated_at
		FROM recipes r
		INNER JOIN recipe_collection_members rcm ON rcm.recipe_id = r.id
		WHERE rcm.collection_id = $1 AND r.user_id = $2
		ORDER BY rcm.added_at DESC
	`
	err := r.db.Select(&recipes, query, collectionID, userID)
	return recipes, err
}

// Exists checks if a recipe with the given name already exists for the user
func (r *RecipeRepository) Exists(userID, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipes WHERE user_id = $1 AND name = $2)`
	err := r.db.Get(&exists, query, userID, name)
	return exists, err
}

// ExistsExcluding checks if a recipe name exists excluding a specific recipe ID
func (r *RecipeRepository) ExistsExcluding(userID, name, excludeID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipes WHERE user_id = $1 AND name = $2 AND id != $3)`
	err := r.db.Get(&exists, query, userID, name, excludeID)
	return exists, err
}
