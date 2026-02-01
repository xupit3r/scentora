package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRecipeVersionRepository_Create(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	version := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Initial version",
	}

	err = versionRepo.Create(version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	if version.ID == "" {
		t.Error("Expected version ID to be set")
	}
	if version.VersionNumber != 1 {
		t.Errorf("Expected version number 1, got %d", version.VersionNumber)
	}
	if !version.IsActive {
		t.Error("Expected version to be active")
	}

	// Verify recipe's active_version_id was updated
	updatedRecipe, err := recipeRepo.FindByID(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find recipe: %v", err)
	}
	if updatedRecipe.ActiveVersionID == nil || *updatedRecipe.ActiveVersionID != version.ID {
		t.Error("Expected recipe's active_version_id to be set")
	}
}

func TestRecipeVersionRepository_CreateMultipleVersions(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create first version
	version1 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 1",
	}
	err = versionRepo.Create(version1)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	// Create second version
	version2 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 2",
	}
	err = versionRepo.Create(version2)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Version 2 should be active, version 1 should be inactive
	if version2.VersionNumber != 2 {
		t.Errorf("Expected version number 2, got %d", version2.VersionNumber)
	}
	if !version2.IsActive {
		t.Error("Expected version 2 to be active")
	}

	// Check version 1 is now inactive
	v1, err := versionRepo.FindByID(version1.ID)
	if err != nil {
		t.Fatalf("Failed to find version 1: %v", err)
	}
	if v1.IsActive {
		t.Error("Expected version 1 to be inactive")
	}
}

func TestRecipeVersionRepository_FindByID(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	notes := "Initial testing notes"
	version := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Test Version",
		Notes:    &notes,
	}
	err = versionRepo.Create(version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	found, err := versionRepo.FindByID(version.ID)
	if err != nil {
		t.Fatalf("Failed to find version: %v", err)
	}

	if found.ID != version.ID {
		t.Errorf("Expected ID %s, got %s", version.ID, found.ID)
	}
	if found.Name != version.Name {
		t.Errorf("Expected name %s, got %s", version.Name, found.Name)
	}
	if found.Notes == nil || *found.Notes != notes {
		t.Error("Expected notes to match")
	}
}

func TestRecipeVersionRepository_FindByRecipeID(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create multiple versions
	for i := 1; i <= 3; i++ {
		version := &models.RecipeVersion{
			RecipeID: recipe.ID,
			Name:     testutil.UniqueString("Version"),
		}
		err = versionRepo.Create(version)
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	versions, err := versionRepo.FindByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find versions: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("Expected 3 versions, got %d", len(versions))
	}

	// Should be ordered by version_number DESC
	for i := 0; i < len(versions)-1; i++ {
		if versions[i].VersionNumber < versions[i+1].VersionNumber {
			t.Error("Expected versions to be ordered descending by version_number")
		}
	}
}

func TestRecipeVersionRepository_FindActiveByRecipeID(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create multiple versions
	version1 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 1",
	}
	err = versionRepo.Create(version1)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	version2 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 2",
	}
	err = versionRepo.Create(version2)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	active, err := versionRepo.FindActiveByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find active version: %v", err)
	}

	if active.ID != version2.ID {
		t.Error("Expected version 2 to be active")
	}
	if !active.IsActive {
		t.Error("Expected active version to have IsActive = true")
	}
}

func TestRecipeVersionRepository_SetActive(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create two versions
	version1 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 1",
	}
	err = versionRepo.Create(version1)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	version2 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 2",
	}
	err = versionRepo.Create(version2)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Version 2 is active, now set version 1 as active
	err = versionRepo.SetActive(version1.ID, recipe.ID)
	if err != nil {
		t.Fatalf("Failed to set active version: %v", err)
	}

	// Verify version 1 is now active
	v1, err := versionRepo.FindByID(version1.ID)
	if err != nil {
		t.Fatalf("Failed to find version 1: %v", err)
	}
	if !v1.IsActive {
		t.Error("Expected version 1 to be active")
	}

	// Verify version 2 is now inactive
	v2, err := versionRepo.FindByID(version2.ID)
	if err != nil {
		t.Fatalf("Failed to find version 2: %v", err)
	}
	if v2.IsActive {
		t.Error("Expected version 2 to be inactive")
	}

	// Verify recipe's active_version_id was updated
	updatedRecipe, err := recipeRepo.FindByID(recipe.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find recipe: %v", err)
	}
	if updatedRecipe.ActiveVersionID == nil || *updatedRecipe.ActiveVersionID != version1.ID {
		t.Error("Expected recipe's active_version_id to be version 1")
	}
}

func TestRecipeVersionRepository_CountByRecipeID(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create 3 versions
	for i := 1; i <= 3; i++ {
		version := &models.RecipeVersion{
			RecipeID: recipe.ID,
			Name:     testutil.UniqueString("Version"),
		}
		err = versionRepo.Create(version)
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	count, err := versionRepo.CountByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to count versions: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestRecipeVersionRepository_Delete(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	
	// Create two versions
	version1 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 1",
	}
	err = versionRepo.Create(version1)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	version2 := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Version 2",
	}
	err = versionRepo.Create(version2)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Delete version 2 (active version)
	err = versionRepo.Delete(version2.ID, recipe.ID)
	if err != nil {
		t.Fatalf("Failed to delete version 2: %v", err)
	}

	// Version 1 should now be active
	v1, err := versionRepo.FindByID(version1.ID)
	if err != nil {
		t.Fatalf("Failed to find version 1: %v", err)
	}
	if !v1.IsActive {
		t.Error("Expected version 1 to be active after deleting version 2")
	}

	// Version 2 should not exist
	_, err = versionRepo.FindByID(version2.ID)
	if err == nil {
		t.Error("Expected error finding deleted version")
	}
}

func TestRecipeVersionRepository_Delete_OnlyVersion(t *testing.T) {
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

	versionRepo := NewRecipeVersionRepository(tdb.DB)
	version := &models.RecipeVersion{
		RecipeID: recipe.ID,
		Name:     "Only Version",
	}
	err = versionRepo.Create(version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Try to delete the only version (should fail)
	err = versionRepo.Delete(version.ID, recipe.ID)
	if err == nil {
		t.Error("Expected error deleting the only version")
	}
}
