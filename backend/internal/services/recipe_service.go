package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeService struct {
	recipeRepo     *repository.RecipeRepository
	versionRepo    *repository.RecipeVersionRepository
	ingredientRepo *repository.RecipeIngredientRepository
	noteRepo       *repository.RecipeNoteRepository
	tagRepo        *repository.RecipeTagRepository
	accordRepo     *repository.AccordRepository
}

func NewRecipeService(
	recipeRepo *repository.RecipeRepository,
	versionRepo *repository.RecipeVersionRepository,
	ingredientRepo *repository.RecipeIngredientRepository,
	noteRepo *repository.RecipeNoteRepository,
	tagRepo *repository.RecipeTagRepository,
	accordRepo *repository.AccordRepository,
) *RecipeService {
	return &RecipeService{
		recipeRepo:     recipeRepo,
		versionRepo:    versionRepo,
		ingredientRepo: ingredientRepo,
		noteRepo:       noteRepo,
		tagRepo:        tagRepo,
		accordRepo:     accordRepo,
	}
}

// CreateRecipe creates a new recipe with validation
func (s *RecipeService) CreateRecipe(userID string, req *models.CreateRecipeRequest) (*models.Recipe, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("recipe name is required")
	}

	// Check for duplicate name
	exists, err := s.recipeRepo.Exists(userID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check recipe existence: %w", err)
	}
	if exists {
		return nil, errors.New("recipe with this name already exists")
	}

	// Set default status if not provided
	status := "draft"
	if req.Status != nil {
		status = *req.Status
	}

	// Create recipe
	recipe := &models.Recipe{
		UserID:         userID,
		Name:           req.Name,
		Description:    req.Description,
		TargetVolumeMl: req.TargetVolumeMl,
		Status:         status,
	}
	err = s.recipeRepo.Create(recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	return recipe, nil
}

// GetRecipe retrieves a recipe by ID with all related data
func (s *RecipeService) GetRecipe(recipeID, userID string) (*models.RecipeResponse, error) {
	// Get recipe
	recipe, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Get tags
	tags, err := s.tagRepo.FindByRecipeID(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}

	// Get active version
	var activeVersion *models.RecipeVersionResponse
	if recipe.ActiveVersionID != nil {
		version, err := s.versionRepo.FindByID(*recipe.ActiveVersionID)
		if err == nil {
			activeVersion = &models.RecipeVersionResponse{
				ID:            version.ID,
				RecipeID:      version.RecipeID,
				VersionNumber: version.VersionNumber,
				Name:          version.Name,
				Notes:         version.Notes,
				IsActive:      version.IsActive,
				CreatedAt:     version.CreatedAt,
			}
		}
	}

	// Build response
	response := &models.RecipeResponse{
		ID:              recipe.ID,
		UserID:          recipe.UserID,
		Name:            recipe.Name,
		Description:     recipe.Description,
		TargetVolumeMl:  recipe.TargetVolumeMl,
		Status:          recipe.Status,
		ActiveVersionID: recipe.ActiveVersionID,
		ActiveVersion:   activeVersion,
		Tags:            tags,
		CreatedAt:       recipe.CreatedAt,
		UpdatedAt:       recipe.UpdatedAt,
	}

	return response, nil
}

// ListRecipes retrieves recipes with filters
func (s *RecipeService) ListRecipes(userID string, status *string, limit, offset int) ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	var err error

	if status != nil && *status != "" {
		recipes, err = s.recipeRepo.FindByStatus(userID, *status)
	} else {
		recipes, err = s.recipeRepo.FindAll(userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list recipes: %w", err)
	}

	return recipes, nil
}

// SearchRecipes searches recipes by query
func (s *RecipeService) SearchRecipes(userID, query string, limit, offset int) ([]*models.Recipe, error) {
	if query == "" {
		return s.ListRecipes(userID, nil, limit, offset)
	}

	recipes, err := s.recipeRepo.Search(userID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}

	return recipes, nil
}

// UpdateRecipe updates a recipe
func (s *RecipeService) UpdateRecipe(recipeID, userID string, req *models.UpdateRecipeRequest) (*models.Recipe, error) {
	// Get existing recipe
	recipe, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Validate name change
	if req.Name != nil && *req.Name != recipe.Name {
		exists, err := s.recipeRepo.ExistsExcluding(userID, *req.Name, recipeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check recipe existence: %w", err)
		}
		if exists {
			return nil, errors.New("recipe with this name already exists")
		}
		recipe.Name = *req.Name
	}

	// Update fields
	if req.Description != nil {
		recipe.Description = req.Description
	}
	if req.Status != nil {
		// Validate status
		validStatuses := map[string]bool{
			"draft": true, "in_progress": true, "tested": true, "finalized": true, "archived": true,
		}
		if !validStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
		recipe.Status = *req.Status
	}

	// Save changes
	err = s.recipeRepo.Update(recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

// DeleteRecipe deletes a recipe and all related data
func (s *RecipeService) DeleteRecipe(recipeID, userID string) error {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Delete recipe (cascade will handle versions, ingredients, notes, tags)
	err = s.recipeRepo.Delete(recipeID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	return nil
}

// GetRecipesByTags retrieves recipes filtered by tags
func (s *RecipeService) GetRecipesByTags(userID string, tags []string, limit, offset int) ([]*models.Recipe, error) {
	if len(tags) == 0 {
		return s.ListRecipes(userID, nil, limit, offset)
	}

	recipes, err := s.recipeRepo.FindByTags(userID, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by tags: %w", err)
	}

	return recipes, nil
}

// GetRecipesByCollection retrieves recipes in a collection
func (s *RecipeService) GetRecipesByCollection(userID, collectionID string) ([]*models.Recipe, error) {
	recipes, err := s.recipeRepo.FindByCollection(collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by collection: %w", err)
	}

	return recipes, nil
}

// GetRecipeStats returns statistics for user's recipes
func (s *RecipeService) GetRecipeStats(userID string) (*models.RecipeStats, error) {
	stats := &models.RecipeStats{}

	// Count by status - repo returns map[string]int
	counts, err := s.recipeRepo.CountByStatus(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count recipes: %w", err)
	}

	stats.DraftCount = counts["draft"]
	stats.InProgressCount = counts["in_progress"]
	stats.TestedCount = counts["tested"]
	stats.FinalizedCount = counts["finalized"]
	stats.ArchivedCount = counts["archived"]
	stats.TotalCount = stats.DraftCount + stats.InProgressCount + stats.TestedCount + stats.FinalizedCount + stats.ArchivedCount

	return stats, nil
}
