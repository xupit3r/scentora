package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeTagRepository struct {
	db *sqlx.DB
}

func NewRecipeTagRepository(db *sqlx.DB) *RecipeTagRepository {
	return &RecipeTagRepository{db: db}
}

// Create adds a tag to a recipe
func (r *RecipeTagRepository) Create(tag *models.RecipeTag) error {
	query := `
		INSERT INTO recipe_tags (recipe_id, tag, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, tag.RecipeID, tag.Tag).Scan(&tag.ID, &tag.CreatedAt)
}

// FindByRecipeID retrieves all tags for a recipe
func (r *RecipeTagRepository) FindByRecipeID(recipeID string) ([]string, error) {
	var tags []string
	query := `
		SELECT tag
		FROM recipe_tags
		WHERE recipe_id = $1
		ORDER BY tag ASC
	`
	err := r.db.Select(&tags, query, recipeID)
	return tags, err
}

// Delete removes a tag from a recipe
func (r *RecipeTagRepository) Delete(recipeID, tag string) error {
	query := `DELETE FROM recipe_tags WHERE recipe_id = $1 AND tag = $2`
	result, err := r.db.Exec(query, recipeID, tag)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}
	return nil
}

// DeleteAll removes all tags from a recipe
func (r *RecipeTagRepository) DeleteAll(recipeID string) error {
	query := `DELETE FROM recipe_tags WHERE recipe_id = $1`
	_, err := r.db.Exec(query, recipeID)
	return err
}

// Exists checks if a tag exists on a recipe
func (r *RecipeTagRepository) Exists(recipeID, tag string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipe_tags WHERE recipe_id = $1 AND tag = $2)`
	err := r.db.Get(&exists, query, recipeID, tag)
	return exists, err
}

// GetPopularTags returns the most commonly used recipe tags with counts
func (r *RecipeTagRepository) GetPopularTags(userID string, limit int) ([]struct {
	Tag   string `db:"tag"`
	Count int    `db:"count"`
}, error) {
	var tags []struct {
		Tag   string `db:"tag"`
		Count int    `db:"count"`
	}
	query := `
		SELECT rt.tag, COUNT(*) as count
		FROM recipe_tags rt
		INNER JOIN recipes r ON r.id = rt.recipe_id
		WHERE r.user_id = $1
		GROUP BY rt.tag
		ORDER BY count DESC, rt.tag ASC
		LIMIT $2
	`
	err := r.db.Select(&tags, query, userID, limit)
	return tags, err
}

// FindRecipesByTag retrieves all recipe IDs that have a specific tag
func (r *RecipeTagRepository) FindRecipesByTag(tag string) ([]string, error) {
	var recipeIDs []string
	query := `
		SELECT recipe_id
		FROM recipe_tags
		WHERE tag = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&recipeIDs, query, tag)
	return recipeIDs, err
}

// CountByRecipeID returns the number of tags for a recipe
func (r *RecipeTagRepository) CountByRecipeID(recipeID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM recipe_tags WHERE recipe_id = $1`
	err := r.db.Get(&count, query, recipeID)
	return count, err
}
