package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func createTestRecipeWithVersion(t *testing.T, tdb *testutil.TestDB) (*models.User, *models.Recipe, *models.RecipeVersion) {
	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        testutil.UniqueEmail("test"),
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	recipeRepo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           testutil.UniqueString("Recipe"),
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = recipeRepo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	version := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 1",
	}
	err = versionRepo.Create(version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	return user, recipe, version
}

func createTestAccord(t *testing.T, tdb *testutil.TestDB, userID string) *models.Accord {
	accordRepo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          userID,
		Name:            testutil.UniqueString("Accord"),
		PyramidPosition: "top",
		VolumeMl:        50.0,
	}
	err := accordRepo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}
	return accord
}

func TestRecipeIngredientRepository_Create(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	accord := createTestAccord(t, tdb, user.ID)

	repo := NewRecipeIngredientRepository(tdb.DB)
	percentage := 25.0
	notes := "Primary note"
	ingredient := &models.RecipeIngredient{
		VersionID:  version.ID,
		AccordID:   accord.ID,
		QuantityMl: 25.0,
		Percentage: &percentage,
		Notes:      &notes,
	}

	err := repo.Create(ingredient)
	if err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	if ingredient.ID == "" {
		t.Error("Expected ingredient ID to be set")
	}
	if ingredient.QuantityDrops != 500 {
		t.Errorf("Expected quantity_drops to be 500 (25.0 * 20), got %d", ingredient.QuantityDrops)
	}
}

func TestRecipeIngredientRepository_FindByVersionID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	
	repo := NewRecipeIngredientRepository(tdb.DB)
	
	// Create multiple ingredients
	for i := 0; i < 3; i++ {
		accord := createTestAccord(t, tdb, user.ID)
		ingredient := &models.RecipeIngredient{
			VersionID:  version.ID,
			AccordID:   accord.ID,
			QuantityMl: 10.0,
		}
		err := repo.Create(ingredient)
		if err != nil {
			t.Fatalf("Failed to create ingredient: %v", err)
		}
	}

	ingredients, err := repo.FindByVersionID(version.ID)
	if err != nil {
		t.Fatalf("Failed to find ingredients: %v", err)
	}

	if len(ingredients) != 3 {
		t.Errorf("Expected 3 ingredients, got %d", len(ingredients))
	}
}

func TestRecipeIngredientRepository_Update(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	accord := createTestAccord(t, tdb, user.ID)

	repo := NewRecipeIngredientRepository(tdb.DB)
	ingredient := &models.RecipeIngredient{
		VersionID:  version.ID,
		AccordID:   accord.ID,
		QuantityMl: 10.0,
	}
	err := repo.Create(ingredient)
	if err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	// Update quantity
	ingredient.QuantityMl = 20.0
	newPercentage := 50.0
	ingredient.Percentage = &newPercentage

	err = repo.Update(ingredient)
	if err != nil {
		t.Fatalf("Failed to update ingredient: %v", err)
	}

	// Verify update
	found, err := repo.FindByID(ingredient.ID)
	if err != nil {
		t.Fatalf("Failed to find ingredient: %v", err)
	}

	if found.QuantityMl != 20.0 {
		t.Errorf("Expected quantity 20.0, got %f", found.QuantityMl)
	}
	if found.QuantityDrops != 400 {
		t.Errorf("Expected quantity_drops 400, got %d", found.QuantityDrops)
	}
}

func TestRecipeIngredientRepository_Delete(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	accord := createTestAccord(t, tdb, user.ID)

	repo := NewRecipeIngredientRepository(tdb.DB)
	ingredient := &models.RecipeIngredient{
		VersionID:  version.ID,
		AccordID:   accord.ID,
		QuantityMl: 10.0,
	}
	err := repo.Create(ingredient)
	if err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	err = repo.Delete(ingredient.ID)
	if err != nil {
		t.Fatalf("Failed to delete ingredient: %v", err)
	}

	_, err = repo.FindByID(ingredient.ID)
	if err == nil {
		t.Error("Expected error finding deleted ingredient")
	}
}

func TestRecipeIngredientRepository_ExistsInVersion(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	accord := createTestAccord(t, tdb, user.ID)

	repo := NewRecipeIngredientRepository(tdb.DB)
	ingredient := &models.RecipeIngredient{
		VersionID:  version.ID,
		AccordID:   accord.ID,
		QuantityMl: 10.0,
	}
	err := repo.Create(ingredient)
	if err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	exists, err := repo.ExistsInVersion(version.ID, accord.ID)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected ingredient to exist in version")
	}

	exists, err = repo.ExistsInVersion(version.ID, "nonexistent-accord-id")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected ingredient not to exist")
	}
}

func TestRecipeIngredientRepository_GetTotalVolume(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)

	repo := NewRecipeIngredientRepository(tdb.DB)
	
	// Create ingredients with different volumes
	volumes := []float64{10.0, 20.0, 15.5}
	expectedTotal := 45.5

	for _, vol := range volumes {
		accord := createTestAccord(t, tdb, user.ID)
		ingredient := &models.RecipeIngredient{
			VersionID:  version.ID,
			AccordID:   accord.ID,
			QuantityMl: vol,
		}
		err := repo.Create(ingredient)
		if err != nil {
			t.Fatalf("Failed to create ingredient: %v", err)
		}
	}

	total, err := repo.GetTotalVolume(version.ID)
	if err != nil {
		t.Fatalf("Failed to get total volume: %v", err)
	}

	if total != expectedTotal {
		t.Errorf("Expected total volume %f, got %f", expectedTotal, total)
	}
}

func TestRecipeIngredientRepository_FindByAccordID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, _, version := createTestRecipeWithVersion(t, tdb)
	accord := createTestAccord(t, tdb, user.ID)

	repo := NewRecipeIngredientRepository(tdb.DB)
	ingredient := &models.RecipeIngredient{
		VersionID:  version.ID,
		AccordID:   accord.ID,
		QuantityMl: 10.0,
	}
	err := repo.Create(ingredient)
	if err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	versionIDs, err := repo.FindByAccordID(accord.ID)
	if err != nil {
		t.Fatalf("Failed to find by accord ID: %v", err)
	}

	if len(versionIDs) != 1 {
		t.Errorf("Expected 1 version ID, got %d", len(versionIDs))
	}
	if versionIDs[0] != version.ID {
		t.Error("Expected correct version ID")
	}
}
