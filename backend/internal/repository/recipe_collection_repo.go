package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeCollectionRepository struct {
	db *sqlx.DB
}

func NewRecipeCollectionRepository(db *sqlx.DB) *RecipeCollectionRepository {
	return &RecipeCollectionRepository{db: db}
}

// Create creates a new recipe collection
func (r *RecipeCollectionRepository) Create(collection *models.RecipeCollection) error {
	query := `
		INSERT INTO recipe_collections (
			user_id, name, description, created_at, updated_at
		)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		collection.UserID,
		collection.Name,
		collection.Description,
	).Scan(&collection.ID, &collection.CreatedAt, &collection.UpdatedAt)
}

// FindByID retrieves a collection by ID
func (r *RecipeCollectionRepository) FindByID(id, userID string) (*models.RecipeCollection, error) {
	var collection models.RecipeCollection
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM recipe_collections
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&collection, query, id, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("collection not found")
	}
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// FindAll retrieves all collections for a user
func (r *RecipeCollectionRepository) FindAll(userID string) ([]*models.RecipeCollection, error) {
	var collections []*models.RecipeCollection
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM recipe_collections
		WHERE user_id = $1
		ORDER BY name ASC
	`
	err := r.db.Select(&collections, query, userID)
	return collections, err
}

// Update updates a collection's name and description
func (r *RecipeCollectionRepository) Update(collection *models.RecipeCollection) error {
	query := `
		UPDATE recipe_collections
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		query,
		collection.Name,
		collection.Description,
		collection.ID,
		collection.UserID,
	).Scan(&collection.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("collection not found")
	}
	return err
}

// Delete deletes a collection (cascade will remove members)
func (r *RecipeCollectionRepository) Delete(id, userID string) error {
	query := `DELETE FROM recipe_collections WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("collection not found")
	}
	return nil
}

// AddRecipe adds a recipe to a collection
func (r *RecipeCollectionRepository) AddRecipe(member *models.RecipeCollectionMember) error {
	query := `
		INSERT INTO recipe_collection_members (
			collection_id, recipe_id, added_at
		)
		VALUES ($1, $2, NOW())
		RETURNING id, added_at
	`
	return r.db.QueryRow(
		query,
		member.CollectionID,
		member.RecipeID,
	).Scan(&member.ID, &member.AddedAt)
}

// RemoveRecipe removes a recipe from a collection
func (r *RecipeCollectionRepository) RemoveRecipe(collectionID, recipeID string) error {
	query := `
		DELETE FROM recipe_collection_members
		WHERE collection_id = $1 AND recipe_id = $2
	`
	result, err := r.db.Exec(query, collectionID, recipeID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("recipe not in collection")
	}
	return nil
}

// GetRecipeIDs retrieves all recipe IDs in a collection
func (r *RecipeCollectionRepository) GetRecipeIDs(collectionID string) ([]string, error) {
	var recipeIDs []string
	query := `
		SELECT recipe_id
		FROM recipe_collection_members
		WHERE collection_id = $1
		ORDER BY added_at DESC
	`
	err := r.db.Select(&recipeIDs, query, collectionID)
	return recipeIDs, err
}

// IsRecipeInCollection checks if a recipe is in a collection
func (r *RecipeCollectionRepository) IsRecipeInCollection(collectionID, recipeID string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM recipe_collection_members 
			WHERE collection_id = $1 AND recipe_id = $2
		)
	`
	err := r.db.Get(&exists, query, collectionID, recipeID)
	return exists, err
}

// CountRecipes returns the number of recipes in a collection
func (r *RecipeCollectionRepository) CountRecipes(collectionID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM recipe_collection_members
		WHERE collection_id = $1
	`
	err := r.db.Get(&count, query, collectionID)
	return count, err
}

// GetCollectionsByRecipeID returns all collections containing a recipe
func (r *RecipeCollectionRepository) GetCollectionsByRecipeID(recipeID, userID string) ([]*models.RecipeCollection, error) {
	var collections []*models.RecipeCollection
	query := `
		SELECT c.id, c.user_id, c.name, c.description, c.created_at, c.updated_at
		FROM recipe_collections c
		INNER JOIN recipe_collection_members m ON m.collection_id = c.id
		WHERE m.recipe_id = $1 AND c.user_id = $2
		ORDER BY c.name ASC
	`
	err := r.db.Select(&collections, query, recipeID, userID)
	return collections, err
}

// Exists checks if a collection name already exists for the user
func (r *RecipeCollectionRepository) Exists(userID, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipe_collections WHERE user_id = $1 AND name = $2)`
	err := r.db.Get(&exists, query, userID, name)
	return exists, err
}

// ExistsExcluding checks if a collection name exists excluding a specific collection ID
func (r *RecipeCollectionRepository) ExistsExcluding(userID, name, excludeID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recipe_collections WHERE user_id = $1 AND name = $2 AND id != $3)`
	err := r.db.Get(&exists, query, userID, name, excludeID)
	return exists, err
}
