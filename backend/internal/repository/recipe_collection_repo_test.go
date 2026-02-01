package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRecipeCollectionRepository_Create(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	desc := "My favorite summer recipes"
	collection := &models.RecipeCollection{
		UserID:      user.ID,
		Name:        "Summer Collection",
		Description: &desc,
	}

	err = repo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	if collection.ID == "" {
		t.Error("Expected collection ID to be set")
	}
}

func TestRecipeCollectionRepository_FindByID(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Test Collection",
	}
	err = repo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	found, err := repo.FindByID(collection.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find collection: %v", err)
	}

	if found.Name != "Test Collection" {
		t.Error("Expected correct collection name")
	}
}

func TestRecipeCollectionRepository_FindAll(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	
	names := []string{"Collection C", "Collection A", "Collection B"}
	for _, name := range names {
		collection := &models.RecipeCollection{
			UserID: user.ID,
			Name:   name,
		}
		err = repo.Create(collection)
		if err != nil {
			t.Fatalf("Failed to create collection: %v", err)
		}
	}

	collections, err := repo.FindAll(user.ID)
	if err != nil {
		t.Fatalf("Failed to find collections: %v", err)
	}

	if len(collections) != 3 {
		t.Errorf("Expected 3 collections, got %d", len(collections))
	}

	// Should be ordered alphabetically
	if collections[0].Name != "Collection A" {
		t.Error("Expected collections to be ordered alphabetically")
	}
}

func TestRecipeCollectionRepository_Update(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Original Name",
	}
	err = repo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	newDesc := "Updated description"
	collection.Name = "Updated Name"
	collection.Description = &newDesc

	err = repo.Update(collection)
	if err != nil {
		t.Fatalf("Failed to update collection: %v", err)
	}

	found, err := repo.FindByID(collection.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find collection: %v", err)
	}

	if found.Name != "Updated Name" {
		t.Error("Expected updated name")
	}
	if found.Description == nil || *found.Description != newDesc {
		t.Error("Expected updated description")
	}
}

func TestRecipeCollectionRepository_Delete(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "To Delete",
	}
	err = repo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	err = repo.Delete(collection.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete collection: %v", err)
	}

	_, err = repo.FindByID(collection.ID, user.ID)
	if err == nil {
		t.Error("Expected error finding deleted collection")
	}
}

func TestRecipeCollectionRepository_AddRecipe(t *testing.T) {
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
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = recipeRepo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	collectionRepo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Test Collection",
	}
	err = collectionRepo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	member := &models.RecipeCollectionMember{
		CollectionID: collection.ID,
		RecipeID:     recipe.ID,
	}
	err = collectionRepo.AddRecipe(member)
	if err != nil {
		t.Fatalf("Failed to add recipe to collection: %v", err)
	}

	if member.ID == "" {
		t.Error("Expected member ID to be set")
	}
}

func TestRecipeCollectionRepository_RemoveRecipe(t *testing.T) {
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
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = recipeRepo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	collectionRepo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Test Collection",
	}
	err = collectionRepo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	member := &models.RecipeCollectionMember{
		CollectionID: collection.ID,
		RecipeID:     recipe.ID,
	}
	err = collectionRepo.AddRecipe(member)
	if err != nil {
		t.Fatalf("Failed to add recipe: %v", err)
	}

	err = collectionRepo.RemoveRecipe(collection.ID, recipe.ID)
	if err != nil {
		t.Fatalf("Failed to remove recipe: %v", err)
	}

	count, err := collectionRepo.CountRecipes(collection.ID)
	if err != nil {
		t.Fatalf("Failed to count recipes: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 recipes after removal, got %d", count)
	}
}

func TestRecipeCollectionRepository_GetRecipeIDs(t *testing.T) {
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
	collectionRepo := NewRecipeCollectionRepository(tdb.DB)
	
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Test Collection",
	}
	err = collectionRepo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Create and add 3 recipes
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

		member := &models.RecipeCollectionMember{
			CollectionID: collection.ID,
			RecipeID:     recipe.ID,
		}
		err = collectionRepo.AddRecipe(member)
		if err != nil {
			t.Fatalf("Failed to add recipe: %v", err)
		}
	}

	recipeIDs, err := collectionRepo.GetRecipeIDs(collection.ID)
	if err != nil {
		t.Fatalf("Failed to get recipe IDs: %v", err)
	}

	if len(recipeIDs) != 3 {
		t.Errorf("Expected 3 recipe IDs, got %d", len(recipeIDs))
	}
}

func TestRecipeCollectionRepository_IsRecipeInCollection(t *testing.T) {
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
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = recipeRepo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	collectionRepo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Test Collection",
	}
	err = collectionRepo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	member := &models.RecipeCollectionMember{
		CollectionID: collection.ID,
		RecipeID:     recipe.ID,
	}
	err = collectionRepo.AddRecipe(member)
	if err != nil {
		t.Fatalf("Failed to add recipe: %v", err)
	}

	exists, err := collectionRepo.IsRecipeInCollection(collection.ID, recipe.ID)
	if err != nil {
		t.Fatalf("Failed to check if recipe in collection: %v", err)
	}
	if !exists {
		t.Error("Expected recipe to be in collection")
	}

	nonexistentID := "00000000-0000-0000-0000-000000000000"
	exists, err = collectionRepo.IsRecipeInCollection(collection.ID, nonexistentID)
	if err != nil {
		t.Fatalf("Failed to check if recipe in collection: %v", err)
	}
	if exists {
		t.Error("Expected recipe not to be in collection")
	}
}

func TestRecipeCollectionRepository_GetCollectionsByRecipeID(t *testing.T) {
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
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Test Recipe",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = recipeRepo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	collectionRepo := NewRecipeCollectionRepository(tdb.DB)
	
	// Create 2 collections and add recipe to both
	for i := 0; i < 2; i++ {
		collection := &models.RecipeCollection{
			UserID: user.ID,
			Name:   testutil.UniqueString("Collection"),
		}
		err = collectionRepo.Create(collection)
		if err != nil {
			t.Fatalf("Failed to create collection: %v", err)
		}

		member := &models.RecipeCollectionMember{
			CollectionID: collection.ID,
			RecipeID:     recipe.ID,
		}
		err = collectionRepo.AddRecipe(member)
		if err != nil {
			t.Fatalf("Failed to add recipe: %v", err)
		}
	}

	collections, err := collectionRepo.GetCollectionsByRecipeID(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to get collections by recipe ID: %v", err)
	}

	if len(collections) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(collections))
	}
}

func TestRecipeCollectionRepository_Exists(t *testing.T) {
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

	repo := NewRecipeCollectionRepository(tdb.DB)
	collection := &models.RecipeCollection{
		UserID: user.ID,
		Name:   "Unique Collection",
	}
	err = repo.Create(collection)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	exists, err := repo.Exists(user.ID, "Unique Collection")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected collection to exist")
	}

	exists, err = repo.Exists(user.ID, "Nonexistent Collection")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected collection not to exist")
	}
}
