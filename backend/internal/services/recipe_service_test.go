package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func setupRecipeService(t *testing.T) (*RecipeService, *testutil.TestDB, string) {
	tdb := testutil.SetupTestDB(t)
	
	// Create repos
	recipeRepo := repository.NewRecipeRepository(tdb.DB)
	versionRepo := repository.NewRecipeVersionRepository(tdb.DB)
	ingredientRepo := repository.NewRecipeIngredientRepository(tdb.DB)
	noteRepo := repository.NewRecipeNoteRepository(tdb.DB)
	tagRepo := repository.NewRecipeTagRepository(tdb.DB)
	accordRepo := repository.NewAccordRepository(tdb.DB)
	userRepo := repository.NewUserRepository(tdb.DB)
	
	// Create test user
	user := &models.User{
		Email:        testutil.UniqueEmail("recipe"),
		Username:     testutil.UniqueString("recipeuser"),
		PasswordHash: "hashed",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)
	
	service := NewRecipeService(recipeRepo, versionRepo, ingredientRepo, noteRepo, tagRepo, accordRepo)
	return service, tdb, user.ID
}

// Test CreateRecipe
func TestRecipeService_CreateRecipe(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateRecipeRequest{
		Name:           "Summer Citrus Blend",
		Description:    stringPtr("A refreshing citrus blend"),
		TargetVolumeMl: 100.0,
	}

	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)
	assert.NotEmpty(t, recipe.ID)
	assert.Equal(t, "Summer Citrus Blend", recipe.Name)
	assert.Equal(t, 100.0, recipe.TargetVolumeMl)
	assert.Equal(t, "draft", recipe.Status)
}

func TestRecipeService_CreateRecipe_WithStatus(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	status := "in_progress"
	req := &models.CreateRecipeRequest{
		Name:           "Test Recipe",
		TargetVolumeMl: 50.0,
		Status:         &status,
	}

	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)
	assert.Equal(t, "in_progress", recipe.Status)
}

func TestRecipeService_CreateRecipe_MissingName(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateRecipeRequest{
		Name:           "",
		TargetVolumeMl: 100.0,
	}

	_, err := service.CreateRecipe(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestRecipeService_CreateRecipe_DuplicateName(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateRecipeRequest{
		Name:           "Duplicate Recipe",
		TargetVolumeMl: 100.0,
	}

	// Create first recipe
	_, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Try to create duplicate
	_, err = service.CreateRecipe(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// Test GetRecipe
func TestRecipeService_GetRecipe(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipe
	req := &models.CreateRecipeRequest{
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
	}
	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Get recipe
	result, err := service.GetRecipe(recipe.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, recipe.ID, result.ID)
	assert.Equal(t, "Test Recipe", result.Name)
}

func TestRecipeService_GetRecipe_NotFound(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	_, err := service.GetRecipe("non-existent-id", userID)
	assert.Error(t, err)
}

func TestRecipeService_GetRecipe_WrongUser(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipe for userID
	req := &models.CreateRecipeRequest{
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
	}
	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Try to access with different user ID
	_, err = service.GetRecipe(recipe.ID, "different-user-id")
	assert.Error(t, err)
}

// Test ListRecipes
func TestRecipeService_ListRecipes(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create some recipes
	for i := 0; i < 3; i++ {
		req := &models.CreateRecipeRequest{
			Name:           testutil.UniqueString("recipe"),
			TargetVolumeMl: 100.0,
		}
		_, err := service.CreateRecipe(userID, req)
		require.NoError(t, err)
	}

	// List recipes
	recipes, err := service.ListRecipes(userID, nil, 100, 0)
	require.NoError(t, err)
	assert.Len(t, recipes, 3)
}

func TestRecipeService_ListRecipes_ByStatus(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipes with different statuses
	status1 := "draft"
	status2 := "tested"
	
	req1 := &models.CreateRecipeRequest{
		Name:           testutil.UniqueString("draft"),
		TargetVolumeMl: 100.0,
		Status:         &status1,
	}
	_, err := service.CreateRecipe(userID, req1)
	require.NoError(t, err)

	req2 := &models.CreateRecipeRequest{
		Name:           testutil.UniqueString("tested"),
		TargetVolumeMl: 100.0,
		Status:         &status2,
	}
	_, err = service.CreateRecipe(userID, req2)
	require.NoError(t, err)

	// List only draft recipes
	recipes, err := service.ListRecipes(userID, &status1, 100, 0)
	require.NoError(t, err)
	assert.Len(t, recipes, 1)
	assert.Equal(t, "draft", recipes[0].Status)
}

// Test UpdateRecipe
func TestRecipeService_UpdateRecipe(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipe
	req := &models.CreateRecipeRequest{
		Name:           "Original Name",
		TargetVolumeMl: 100.0,
	}
	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Update recipe
	updateReq := &models.UpdateRecipeRequest{
		Name:        stringPtr("Updated Name"),
		Description: stringPtr("New description"),
	}
	updated, err := service.UpdateRecipe(recipe.ID, userID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "New description", *updated.Description)
}

func TestRecipeService_UpdateRecipe_Status(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipe
	req := &models.CreateRecipeRequest{
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
	}
	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Update status
	newStatus := "tested"
	updateReq := &models.UpdateRecipeRequest{
		Status: &newStatus,
	}
	updated, err := service.UpdateRecipe(recipe.ID, userID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "tested", updated.Status)
}

// Test DeleteRecipe
func TestRecipeService_DeleteRecipe(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create recipe
	req := &models.CreateRecipeRequest{
		Name:           "To Delete",
		TargetVolumeMl: 100.0,
	}
	recipe, err := service.CreateRecipe(userID, req)
	require.NoError(t, err)

	// Delete recipe
	err = service.DeleteRecipe(recipe.ID, userID)
	require.NoError(t, err)

	// Verify deleted
	_, err = service.GetRecipe(recipe.ID, userID)
	assert.Error(t, err)
}

func TestRecipeService_DeleteRecipe_NotFound(t *testing.T) {
	service, tdb, userID := setupRecipeService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	err := service.DeleteRecipe("non-existent-id", userID)
	assert.Error(t, err)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

func timePtr(t time.Time) *time.Time {
	return &t
}
