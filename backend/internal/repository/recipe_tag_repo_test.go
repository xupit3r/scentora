package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRecipeTagRepository_Create(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeTagRepository(tdb.DB)
	tag := &models.RecipeTag{
		RecipeID: recipe.ID,
		Tag:      "floral",
	}

	err := repo.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	if tag.ID == "" {
		t.Error("Expected tag ID to be set")
	}
}

func TestRecipeTagRepository_FindByRecipeID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeTagRepository(tdb.DB)
	
	tags := []string{"floral", "citrus", "fresh"}
	for _, tagName := range tags {
		tag := &models.RecipeTag{
			RecipeID: recipe.ID,
			Tag:      tagName,
		}
		err := repo.Create(tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}
	}

	foundTags, err := repo.FindByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find tags: %v", err)
	}

	if len(foundTags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(foundTags))
	}

	// Should be ordered alphabetically
	expectedOrder := []string{"citrus", "floral", "fresh"}
	for i, expected := range expectedOrder {
		if foundTags[i] != expected {
			t.Errorf("Expected tag %s at position %d, got %s", expected, i, foundTags[i])
		}
	}
}

func TestRecipeTagRepository_Delete(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeTagRepository(tdb.DB)
	tag := &models.RecipeTag{
		RecipeID: recipe.ID,
		Tag:      "floral",
	}
	err := repo.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	err = repo.Delete(recipe.ID, "floral")
	if err != nil {
		t.Fatalf("Failed to delete tag: %v", err)
	}

	tags, err := repo.FindByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find tags: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags after deletion, got %d", len(tags))
	}
}

func TestRecipeTagRepository_DeleteAll(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeTagRepository(tdb.DB)
	
	tags := []string{"floral", "citrus", "fresh"}
	for _, tagName := range tags {
		tag := &models.RecipeTag{
			RecipeID: recipe.ID,
			Tag:      tagName,
		}
		err := repo.Create(tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}
	}

	err := repo.DeleteAll(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to delete all tags: %v", err)
	}

	foundTags, err := repo.FindByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find tags: %v", err)
	}

	if len(foundTags) != 0 {
		t.Errorf("Expected 0 tags after delete all, got %d", len(foundTags))
	}
}

func TestRecipeTagRepository_Exists(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeTagRepository(tdb.DB)
	tag := &models.RecipeTag{
		RecipeID: recipe.ID,
		Tag:      "floral",
	}
	err := repo.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	exists, err := repo.Exists(recipe.ID, "floral")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected tag to exist")
	}

	exists, err = repo.Exists(recipe.ID, "nonexistent")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected tag not to exist")
	}
}

func TestRecipeTagRepository_GetPopularTags(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

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
	tagRepo := NewRecipeTagRepository(tdb.DB)

	// Create 3 recipes
	recipes := make([]*models.Recipe, 3)
	for i := 0; i < 3; i++ {
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
		recipes[i] = recipe
	}

	// Tag all 3 with "floral", 2 with "citrus", 1 with "woody"
	for _, recipe := range recipes {
		tag := &models.RecipeTag{RecipeID: recipe.ID, Tag: "floral"}
		tagRepo.Create(tag)
	}
	for i := 0; i < 2; i++ {
		tag := &models.RecipeTag{RecipeID: recipes[i].ID, Tag: "citrus"}
		tagRepo.Create(tag)
	}
	tag := &models.RecipeTag{RecipeID: recipes[0].ID, Tag: "woody"}
	tagRepo.Create(tag)

	popularTags, err := tagRepo.GetPopularTags(user.ID, 10)
	if err != nil {
		t.Fatalf("Failed to get popular tags: %v", err)
	}

	if len(popularTags) != 3 {
		t.Errorf("Expected 3 popular tags, got %d", len(popularTags))
	}

	// Should be ordered by count DESC
	if popularTags[0].Tag != "floral" || popularTags[0].Count != 3 {
		t.Error("Expected 'floral' to be most popular with count 3")
	}
	if popularTags[1].Tag != "citrus" || popularTags[1].Count != 2 {
		t.Error("Expected 'citrus' to be second with count 2")
	}
	if popularTags[2].Tag != "woody" || popularTags[2].Count != 1 {
		t.Error("Expected 'woody' to be third with count 1")
	}
}
