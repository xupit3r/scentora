package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRecipeRepository_Create(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Summer Breeze",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}

	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	if recipe.ID == "" {
		t.Error("Expected recipe ID to be set")
	}
	if recipe.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
	if recipe.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}
}

func TestRecipeRepository_CreateWithDescription(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	desc := "A fresh summer fragrance"
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Summer Breeze",
		Description:    &desc,
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}

	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	if recipe.Description == nil || *recipe.Description != desc {
		t.Error("Expected description to be set")
	}
}

func TestRecipeRepository_FindByID(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Summer Breeze",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	found, err := repo.FindByID(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find recipe: %v", err)
	}

	if found.ID != recipe.ID {
		t.Errorf("Expected ID %s, got %s", recipe.ID, found.ID)
	}
	if found.Name != recipe.Name {
		t.Errorf("Expected name %s, got %s", recipe.Name, found.Name)
	}
}

func TestRecipeRepository_FindByID_NotFound(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	_, err = repo.FindByID("nonexistent-id", user.ID)
	if err == nil {
		t.Error("Expected error for nonexistent recipe")
	}
}

func TestRecipeRepository_FindAll(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	
	// Create multiple recipes
	for i := 0; i < 3; i++ {
		recipe := &models.Recipe{
			UserID:         user.ID,
			Name:           testutil.UniqueString("Recipe"),
			TargetVolumeMl: 100.0,
			Status:         "draft",
		}
		err = repo.Create(recipe)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}
	}

	recipes, err := repo.FindAll(user.ID)
	if err != nil {
		t.Fatalf("Failed to find recipes: %v", err)
	}

	if len(recipes) != 3 {
		t.Errorf("Expected 3 recipes, got %d", len(recipes))
	}
}

func TestRecipeRepository_FindByStatus(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	
	// Create recipes with different statuses
	statuses := []string{"draft", "tested", "finalized"}
	for _, status := range statuses {
		recipe := &models.Recipe{
			UserID:         user.ID,
			Name:           testutil.UniqueString("Recipe"),
			TargetVolumeMl: 100.0,
			Status:         status,
		}
		err = repo.Create(recipe)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}
	}

	tested, err := repo.FindByStatus(user.ID, "tested")
	if err != nil {
		t.Fatalf("Failed to find tested recipes: %v", err)
	}

	if len(tested) != 1 {
		t.Errorf("Expected 1 tested recipe, got %d", len(tested))
	}
	if tested[0].Status != "tested" {
		t.Errorf("Expected status 'tested', got '%s'", tested[0].Status)
	}
}

func TestRecipeRepository_Search(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	
	// Create recipes with searchable names
	desc1 := "A citrus summer fragrance"
	recipe1 := &models.Recipe{
		UserID:         user.ID,
		Name:           "Summer Citrus",
		Description:    &desc1,
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe1)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	desc2 := "A woody winter scent"
	recipe2 := &models.Recipe{
		UserID:         user.ID,
		Name:           "Winter Woods",
		Description:    &desc2,
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe2)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Search by name
	results, err := repo.Search(user.ID, "citrus")
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Summer Citrus" {
		t.Errorf("Expected 'Summer Citrus', got '%s'", results[0].Name)
	}

	// Search by description
	results, err = repo.Search(user.ID, "woody")
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestRecipeRepository_Update(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Original Name",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Update recipe
	recipe.Name = "Updated Name"
	recipe.Status = "tested"
	recipe.TargetVolumeMl = 150.0

	err = repo.Update(recipe)
	if err != nil {
		t.Fatalf("Failed to update recipe: %v", err)
	}

	// Verify update
	found, err := repo.FindByID(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find recipe: %v", err)
	}

	if found.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", found.Name)
	}
	if found.Status != "tested" {
		t.Errorf("Expected status 'tested', got '%s'", found.Status)
	}
	if found.TargetVolumeMl != 150.0 {
		t.Errorf("Expected target volume 150.0, got %f", found.TargetVolumeMl)
	}
}

func TestRecipeRepository_Delete(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "To Delete",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	err = repo.Delete(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete recipe: %v", err)
	}

	_, err = repo.FindByID(recipe.ID, user.ID)
	if err == nil {
		t.Error("Expected error finding deleted recipe")
	}
}

func TestRecipeRepository_CountByStatus(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	
	// Create recipes with different statuses
	statuses := map[string]int{
		"draft":     2,
		"tested":    1,
		"finalized": 3,
	}

	for status, count := range statuses {
		for i := 0; i < count; i++ {
			recipe := &models.Recipe{
				UserID:         user.ID,
				Name:           testutil.UniqueString("Recipe"),
				TargetVolumeMl: 100.0,
				Status:         status,
			}
			err = repo.Create(recipe)
			if err != nil {
				t.Fatalf("Failed to create recipe: %v", err)
			}
		}
	}

	counts, err := repo.CountByStatus(user.ID)
	if err != nil {
		t.Fatalf("Failed to count by status: %v", err)
	}

	for status, expectedCount := range statuses {
		if counts[status] != expectedCount {
			t.Errorf("Expected %d recipes with status '%s', got %d", expectedCount, status, counts[status])
		}
	}
}

func TestRecipeRepository_Exists(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe := &models.Recipe{
		UserID:         user.ID,
		Name:           "Unique Name",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	exists, err := repo.Exists(user.ID, "Unique Name")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected recipe to exist")
	}

	exists, err = repo.Exists(user.ID, "Nonexistent Name")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected recipe not to exist")
	}
}

func TestRecipeRepository_ExistsExcluding(t *testing.T) {
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

	repo := NewRecipeRepository(tdb.DB)
	recipe1 := &models.Recipe{
		UserID:         user.ID,
		Name:           "Recipe One",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe1)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	recipe2 := &models.Recipe{
		UserID:         user.ID,
		Name:           "Recipe Two",
		TargetVolumeMl: 100.0,
		Status:         "draft",
	}
	err = repo.Create(recipe2)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Check if Recipe One exists excluding itself (should be false)
	exists, err := repo.ExistsExcluding(user.ID, "Recipe One", recipe1.ID)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected recipe not to exist when excluding itself")
	}

	// Check if Recipe Two exists excluding Recipe One (should be true)
	exists, err = repo.ExistsExcluding(user.ID, "Recipe Two", recipe1.ID)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected recipe to exist")
	}
}
