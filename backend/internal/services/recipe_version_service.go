package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeVersionService struct {
	versionRepo    *repository.RecipeVersionRepository
	recipeRepo     *repository.RecipeRepository
	ingredientRepo *repository.RecipeIngredientRepository
}

func NewRecipeVersionService(
	versionRepo *repository.RecipeVersionRepository,
	recipeRepo *repository.RecipeRepository,
	ingredientRepo *repository.RecipeIngredientRepository,
) *RecipeVersionService {
	return &RecipeVersionService{
		versionRepo:    versionRepo,
		recipeRepo:     recipeRepo,
		ingredientRepo: ingredientRepo,
	}
}

// CreateVersion creates a new version for a recipe
func (s *RecipeVersionService) CreateVersion(recipeID, userID string, req *models.CreateRecipeVersionRequest) (*models.RecipeVersion, error) {
	// Verify recipe exists and belongs to user
	recipe, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Get next version number
	count, err := s.versionRepo.CountByRecipeID(recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count versions: %w", err)
	}

	// Create version
	version := &models.RecipeVersion{
		RecipeID:      recipe.ID,
		VersionNumber: count + 1,
		Name:          fmt.Sprintf("v%d", count+1),
		Notes:         req.Notes,
		IsActive:      false, // Will be set by repo
	}
	
	err = s.versionRepo.Create(version)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	return version, nil
}

// GetVersion retrieves a version by ID with ingredients
func (s *RecipeVersionService) GetVersion(versionID string) (*models.RecipeVersionResponse, error) {
	// Get version
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Get ingredients
	ingredients, err := s.ingredientRepo.FindByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load ingredients: %w", err)
	}

	// Convert ingredients to response format
	ingredientResponses := make([]models.RecipeIngredientResponse, len(ingredients))
	for i, ing := range ingredients {
		ingredientResponses[i] = models.RecipeIngredientResponse{
			ID:            ing.ID,
			VersionID:     ing.VersionID,
			AccordID:      ing.AccordID,
			QuantityMl:    ing.QuantityMl,
			QuantityDrops: ing.QuantityDrops,
			Percentage:    ing.Percentage,
			Notes:         ing.Notes,
			CreatedAt:     ing.CreatedAt,
		}
	}

	response := &models.RecipeVersionResponse{
		ID:            version.ID,
		RecipeID:      version.RecipeID,
		VersionNumber: version.VersionNumber,
		Name:          version.Name,
		Notes:         version.Notes,
		IsActive:      version.IsActive,
		Ingredients:   ingredientResponses,
		CreatedAt:     version.CreatedAt,
	}

	return response, nil
}

// ListVersions retrieves all versions for a recipe
func (s *RecipeVersionService) ListVersions(recipeID, userID string) ([]*models.RecipeVersion, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	versions, err := s.versionRepo.FindByRecipeID(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	return versions, nil
}

// GetActiveVersion retrieves the active version for a recipe
func (s *RecipeVersionService) GetActiveVersion(recipeID, userID string) (*models.RecipeVersionResponse, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Get active version
	version, err := s.versionRepo.FindActiveByRecipeID(recipeID)
	if err != nil {
		return nil, fmt.Errorf("no active version found: %w", err)
	}

	// Get ingredients
	ingredients, err := s.ingredientRepo.FindByVersionID(version.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load ingredients: %w", err)
	}

	// Convert ingredients to response format
	ingredientResponses := make([]models.RecipeIngredientResponse, len(ingredients))
	for i, ing := range ingredients {
		ingredientResponses[i] = models.RecipeIngredientResponse{
			ID:            ing.ID,
			VersionID:     ing.VersionID,
			AccordID:      ing.AccordID,
			QuantityMl:    ing.QuantityMl,
			QuantityDrops: ing.QuantityDrops,
			Percentage:    ing.Percentage,
			Notes:         ing.Notes,
			CreatedAt:     ing.CreatedAt,
		}
	}

	response := &models.RecipeVersionResponse{
		ID:            version.ID,
		RecipeID:      version.RecipeID,
		VersionNumber: version.VersionNumber,
		Name:          version.Name,
		Notes:         version.Notes,
		IsActive:      version.IsActive,
		Ingredients:   ingredientResponses,
		CreatedAt:     version.CreatedAt,
	}

	return response, nil
}

// SetActiveVersion sets the active version for a recipe
func (s *RecipeVersionService) SetActiveVersion(recipeID, versionID, userID string) error {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(recipeID, userID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Verify version belongs to this recipe
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}
	if version.RecipeID != recipeID {
		return errors.New("version does not belong to this recipe")
	}

	// Set active (repo handles deactivating others)
	err = s.versionRepo.SetActive(recipeID, versionID)
	if err != nil {
		return fmt.Errorf("failed to set active version: %w", err)
	}

	return nil
}

// DeleteVersion deletes a version
func (s *RecipeVersionService) DeleteVersion(versionID, userID string) error {
	// Get version
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Delete version (repo handles "only version" protection and auto-activation)
	err = s.versionRepo.Delete(versionID, version.RecipeID)
	if err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}

	return nil
}

// DuplicateVersion creates a copy of a version
func (s *RecipeVersionService) DuplicateVersion(versionID, userID string) (*models.RecipeVersion, error) {
	// Get existing version
	existingVersion, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(existingVersion.RecipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Get ingredients from existing version
	ingredients, err := s.ingredientRepo.FindByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load ingredients: %w", err)
	}

	// Get next version number
	count, err := s.versionRepo.CountByRecipeID(existingVersion.RecipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to count versions: %w", err)
	}

	// Create new version
	newVersion := &models.RecipeVersion{
		RecipeID:      existingVersion.RecipeID,
		VersionNumber: count + 1,
		Name:          fmt.Sprintf("v%d", count+1),
		Notes:         existingVersion.Notes,
		IsActive:      false,
	}
	
	err = s.versionRepo.Create(newVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Copy ingredients to new version
	for _, ingredient := range ingredients {
		newIngredient := &models.RecipeIngredient{
			VersionID:  newVersion.ID,
			AccordID:   ingredient.AccordID,
			QuantityMl: ingredient.QuantityMl,
		}
		err = s.ingredientRepo.Create(newIngredient)
		if err != nil {
			// If ingredient copy fails, clean up and return error
			_ = s.versionRepo.Delete(newVersion.ID, newVersion.RecipeID)
			return nil, fmt.Errorf("failed to copy ingredients: %w", err)
		}
	}

	return newVersion, nil
}
