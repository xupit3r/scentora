package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeIngredientService struct {
	ingredientRepo *repository.RecipeIngredientRepository
	versionRepo    *repository.RecipeVersionRepository
	recipeRepo     *repository.RecipeRepository
	accordRepo     *repository.AccordRepository
	userRepo       *repository.UserRepository
}

func NewRecipeIngredientService(
	ingredientRepo *repository.RecipeIngredientRepository,
	versionRepo *repository.RecipeVersionRepository,
	recipeRepo *repository.RecipeRepository,
	accordRepo *repository.AccordRepository,
	userRepo *repository.UserRepository,
) *RecipeIngredientService {
	return &RecipeIngredientService{
		ingredientRepo: ingredientRepo,
		versionRepo:    versionRepo,
		recipeRepo:     recipeRepo,
		accordRepo:     accordRepo,
		userRepo:       userRepo,
	}
}

// AddIngredient adds an ingredient to a recipe version
func (s *RecipeIngredientService) AddIngredient(versionID, userID string, req *models.CreateRecipeIngredientRequest) (*models.RecipeIngredient, error) {
	// Get version and verify recipe ownership
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Validate quantity
	if req.QuantityMl <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Verify accord exists and belongs to user
	accord, err := s.accordRepo.FindByID(req.AccordID, userID)
	if err != nil {
		return nil, fmt.Errorf("accord not found: %w", err)
	}

	// Check if ingredient already exists in this version
	exists, err := s.ingredientRepo.ExistsInVersion(versionID, req.AccordID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ingredient existence: %w", err)
	}
	if exists {
		return nil, errors.New("this accord is already in this version")
	}

	// Volume validation (if enabled for user)
	user, err := s.userRepo.FindByID(userID)
	if err == nil && user.ValidateRecipeVolumes {
		// Check if accord has enough volume
		if req.QuantityMl > accord.VolumeMl {
			return nil, fmt.Errorf("insufficient accord volume: need %.2f ml, have %.2f ml", req.QuantityMl, accord.VolumeMl)
		}
		// Note: Advanced volume checking across versions would require
		// fetching all ingredients for each version, which is expensive.
		// This basic check ensures single-version doesn't exceed accord volume.
	}

	// Create ingredient
	ingredient := &models.RecipeIngredient{
		VersionID:  versionID,
		AccordID:   req.AccordID,
		QuantityMl: req.QuantityMl,
	}

	err = s.ingredientRepo.Create(ingredient)
	if err != nil {
		return nil, fmt.Errorf("failed to add ingredient: %w", err)
	}

	return ingredient, nil
}

// GetIngredients retrieves all ingredients for a version
func (s *RecipeIngredientService) GetIngredients(versionID, userID string) ([]*models.RecipeIngredient, error) {
	// Get version and verify recipe ownership
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	ingredients, err := s.ingredientRepo.FindByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients: %w", err)
	}

	return ingredients, nil
}

// UpdateIngredient updates an ingredient's quantity
func (s *RecipeIngredientService) UpdateIngredient(ingredientID, userID string, req *models.UpdateRecipeIngredientRequest) (*models.RecipeIngredient, error) {
	// Get ingredient
	ingredient, err := s.ingredientRepo.FindByID(ingredientID)
	if err != nil {
		return nil, fmt.Errorf("ingredient not found: %w", err)
	}

	// Get version and verify recipe ownership
	version, err := s.versionRepo.FindByID(ingredient.VersionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Validate quantity
	if req.QuantityMl <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Volume validation (if enabled for user)
	user, err := s.userRepo.FindByID(userID)
	if err == nil && user.ValidateRecipeVolumes {
		accord, err := s.accordRepo.FindByID(ingredient.AccordID, userID)
		if err != nil {
			return nil, fmt.Errorf("accord not found: %w", err)
		}

		// Check if accord has enough volume for this update
		if req.QuantityMl > accord.VolumeMl {
			return nil, fmt.Errorf("insufficient accord volume: need %.2f ml, have %.2f ml", req.QuantityMl, accord.VolumeMl)
		}
	}

	// Update quantity
	ingredient.QuantityMl = req.QuantityMl
	err = s.ingredientRepo.Update(ingredient)
	if err != nil {
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}

	return ingredient, nil
}

// DeleteIngredient removes an ingredient from a version
func (s *RecipeIngredientService) DeleteIngredient(ingredientID, userID string) error {
	// Get ingredient
	ingredient, err := s.ingredientRepo.FindByID(ingredientID)
	if err != nil {
		return fmt.Errorf("ingredient not found: %w", err)
	}

	// Get version and verify recipe ownership
	version, err := s.versionRepo.FindByID(ingredient.VersionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Delete ingredient
	err = s.ingredientRepo.Delete(ingredientID)
	if err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	return nil
}

// GetTotalVolume calculates the total volume for a version
func (s *RecipeIngredientService) GetTotalVolume(versionID, userID string) (float64, error) {
	// Get version and verify recipe ownership
	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return 0, fmt.Errorf("version not found: %w", err)
	}

	_, err = s.recipeRepo.FindByID(version.RecipeID, userID)
	if err != nil {
		return 0, fmt.Errorf("recipe not found: %w", err)
	}

	totalVolume, err := s.ingredientRepo.GetTotalVolume(versionID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total volume: %w", err)
	}

	return totalVolume, nil
}
