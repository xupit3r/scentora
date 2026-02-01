package services

import (
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeTagService struct {
	tagRepo    *repository.RecipeTagRepository
	recipeRepo *repository.RecipeRepository
}

func NewRecipeTagService(
	tagRepo *repository.RecipeTagRepository,
	recipeRepo *repository.RecipeRepository,
) *RecipeTagService {
	return &RecipeTagService{
		tagRepo:    tagRepo,
		recipeRepo: recipeRepo,
	}
}

// AddTag adds a tag to a recipe
func (s *RecipeTagService) AddTag(recipeID, userID, tag string) error {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Check if tag already exists
	exists, err := s.tagRepo.Exists(recipeID, tag)
	if err != nil {
		return fmt.Errorf("failed to check tag existence: %w", err)
	}
	if exists {
		// Already exists, no error
		return nil
	}

	// Add tag
	recipeTag := &models.RecipeTag{
		RecipeID: recipeID,
		Tag:      tag,
	}
	err = s.tagRepo.Create(recipeTag)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

// RemoveTag removes a tag from a recipe
func (s *RecipeTagService) RemoveTag(recipeID, userID, tag string) error {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Remove tag
	err = s.tagRepo.Delete(recipeID, tag)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

// GetTags retrieves all tags for a recipe
func (s *RecipeTagService) GetTags(recipeID, userID string) ([]string, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	tags, err := s.tagRepo.FindByRecipeID(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}

// GetPopularTags retrieves the most popular tags for a user
func (s *RecipeTagService) GetPopularTags(userID string, limit int) ([]struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}, error) {
	if limit <= 0 {
		limit = 10
	}

	tags, err := s.tagRepo.GetPopularTags(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}

	return tags, nil
}
