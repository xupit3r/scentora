package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeCollectionService struct {
	collectionRepo *repository.RecipeCollectionRepository
	recipeRepo     *repository.RecipeRepository
}

func NewRecipeCollectionService(
	collectionRepo *repository.RecipeCollectionRepository,
	recipeRepo *repository.RecipeRepository,
) *RecipeCollectionService {
	return &RecipeCollectionService{
		collectionRepo: collectionRepo,
		recipeRepo:     recipeRepo,
	}
}

// CreateCollection creates a new collection
func (s *RecipeCollectionService) CreateCollection(userID string, req *models.CreateRecipeCollectionRequest) (*models.RecipeCollection, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("collection name is required")
	}

	// Check for duplicate name
	exists, err := s.collectionRepo.Exists(userID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check collection existence: %w", err)
	}
	if exists {
		return nil, errors.New("collection with this name already exists")
	}

	// Create collection
	collection := &models.RecipeCollection{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}
	err = s.collectionRepo.Create(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

// GetCollection retrieves a collection by ID
func (s *RecipeCollectionService) GetCollection(collectionID, userID string) (*models.RecipeCollection, error) {
	collection, err := s.collectionRepo.FindByID(collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// Get recipe count
	if err != nil {
		return nil, fmt.Errorf("failed to count recipes: %w", err)
	}

	return collection, nil
}

// ListCollections retrieves all collections for a user
func (s *RecipeCollectionService) ListCollections(userID string, limit, offset int) ([]*models.RecipeCollection, error) {
	collections, err := s.collectionRepo.FindAll(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	return collections, nil
}

// UpdateCollection updates a collection
func (s *RecipeCollectionService) UpdateCollection(collectionID, userID string, req *models.UpdateRecipeCollectionRequest) (*models.RecipeCollection, error) {
	// Get existing collection
	collection, err := s.collectionRepo.FindByID(collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// Validate name change
	if req.Name != nil && *req.Name != collection.Name {
		if *req.Name == "" {
			return nil, errors.New("collection name is required")
		}

		exists, err := s.collectionRepo.Exists(userID, *req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check collection existence: %w", err)
		}
		if exists {
			return nil, errors.New("collection with this name already exists")
		}
		collection.Name = *req.Name
	}

	// Update description
	if req.Description != nil {
		collection.Description = req.Description
	}

	// Save changes
	err = s.collectionRepo.Update(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	return collection, nil
}

// DeleteCollection deletes a collection
func (s *RecipeCollectionService) DeleteCollection(collectionID, userID string) error {
	// Verify collection exists and belongs to user
	_, err := s.collectionRepo.FindByID(collectionID, userID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Delete collection (cascade will handle membership records)
	err = s.collectionRepo.Delete(userID, collectionID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// AddRecipeToCollection adds a recipe to a collection
func (s *RecipeCollectionService) AddRecipeToCollection(collectionID, recipeID, userID string) error {
	// Verify collection exists and belongs to user
	_, err := s.collectionRepo.FindByID(collectionID, userID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Check if already in collection
	exists, err := s.collectionRepo.IsRecipeInCollection(collectionID, recipeID)
	if err != nil {
		return fmt.Errorf("failed to check recipe membership: %w", err)
	}
	if exists {
		// Already in collection, no error
		return nil
	}

	// Add recipe to collection
	err = s.collectionRepo.AddRecipe(&models.RecipeCollectionMember{CollectionID: collectionID, RecipeID: recipeID})
	if err != nil {
		return fmt.Errorf("failed to add recipe to collection: %w", err)
	}

	return nil
}

// RemoveRecipeFromCollection removes a recipe from a collection
func (s *RecipeCollectionService) RemoveRecipeFromCollection(collectionID, recipeID, userID string) error {
	// Verify collection exists and belongs to user
	_, err := s.collectionRepo.FindByID(collectionID, userID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Remove recipe from collection
	err = s.collectionRepo.RemoveRecipe(collectionID, recipeID)
	if err != nil {
		return fmt.Errorf("failed to remove recipe from collection: %w", err)
	}

	return nil
}

// GetRecipeCollections retrieves all collections containing a recipe
func (s *RecipeCollectionService) GetRecipeCollections(recipeID, userID string) ([]*models.RecipeCollection, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	collections, err := s.collectionRepo.GetCollectionsByRecipeID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	return collections, nil
}
