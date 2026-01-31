package repository

import (
	"testing"
	"time"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestAccordRepository_Create(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create user
	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create accord
	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}

	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	if accord.ID == "" {
		t.Error("Expected accord ID to be set")
	}

	if accord.VolumeDrops != 510 {
		t.Errorf("Expected volume_drops to be 510 (25.5 * 20), got %d", accord.VolumeDrops)
	}
}

func TestAccordRepository_CreateWithOptionalFields(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	supplier := "Perfumer's Apprentice"
	purchaseDate := time.Now()
	dilution := 10.0
	notes := "Very bright and citrusy"

	accord := &models.Accord{
		UserID:             user.ID,
		Name:               "Citrus Fresh",
		PyramidPosition:    "top",
		VolumeMl:           25.5,
		Supplier:           &supplier,
		PurchaseDate:       &purchaseDate,
		DilutionPercentage: &dilution,
		Notes:              &notes,
	}

	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	if accord.Supplier == nil || *accord.Supplier != supplier {
		t.Errorf("Expected supplier %s, got %v", supplier, accord.Supplier)
	}

	if accord.DilutionPercentage == nil || *accord.DilutionPercentage != dilution {
		t.Errorf("Expected dilution %f, got %v", dilution, accord.DilutionPercentage)
	}
}

func TestAccordRepository_FindByID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	created := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(created)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	found, err := repo.FindByID(created.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find accord: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, found.ID)
	}

	if found.Name != "Citrus Fresh" {
		t.Errorf("Expected name Citrus Fresh, got %s", found.Name)
	}
}

func TestAccordRepository_FindByIDNotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	_, err = repo.FindByID("00000000-0000-0000-0000-000000000000", user.ID)
	if err == nil {
		t.Error("Expected error for nonexistent accord, got nil")
	}
}

func TestAccordRepository_List(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user1 := &models.User{
		Email:        "user1@example.com",
		Username:     "user1",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user1)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2 := &models.User{
		Email:        "user2@example.com",
		Username:     "user2",
		PasswordHash: "hashedpassword",
	}
	err = userRepo.Create(user2)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)

	// Create accords for user1
	accord1 := &models.Accord{
		UserID:          user1.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord1)
	if err != nil {
		t.Fatalf("Failed to create accord1: %v", err)
	}

	accord2 := &models.Accord{
		UserID:          user1.ID,
		Name:            "Woody Base",
		PyramidPosition: "base",
		VolumeMl:        30.0,
	}
	err = repo.Create(accord2)
	if err != nil {
		t.Fatalf("Failed to create accord2: %v", err)
	}

	// Create accord for user2
	accord3 := &models.Accord{
		UserID:          user2.ID,
		Name:            "Floral Heart",
		PyramidPosition: "middle",
		VolumeMl:        20.0,
	}
	err = repo.Create(accord3)
	if err != nil {
		t.Fatalf("Failed to create accord3: %v", err)
	}

	// List accords for user1
	accords, err := repo.List(user1.ID)
	if err != nil {
		t.Fatalf("Failed to list accords: %v", err)
	}

	if len(accords) != 2 {
		t.Errorf("Expected 2 accords for user1, got %d", len(accords))
	}
}

func TestAccordRepository_Update(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	// Update accord
	accord.Name = "Citrus Updated"
	accord.VolumeMl = 30.0

	err = repo.Update(accord)
	if err != nil {
		t.Fatalf("Failed to update accord: %v", err)
	}

	if accord.VolumeDrops != 600 {
		t.Errorf("Expected volume_drops to be 600 (30.0 * 20), got %d", accord.VolumeDrops)
	}

	// Verify update
	found, err := repo.FindByID(accord.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to find updated accord: %v", err)
	}

	if found.Name != "Citrus Updated" {
		t.Errorf("Expected name Citrus Updated, got %s", found.Name)
	}

	if found.VolumeMl != 30.0 {
		t.Errorf("Expected volume_ml 30.0, got %f", found.VolumeMl)
	}
}

func TestAccordRepository_Delete(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	// Delete accord
	err = repo.Delete(accord.ID, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete accord: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByID(accord.ID, user.ID)
	if err == nil {
		t.Error("Expected error for deleted accord, got nil")
	}
}

func TestAccordRepository_AddTag(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	// Add tag
	err = repo.AddTag(accord.ID, "fresh")
	if err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}

	// Get tags
	tags, err := repo.GetTagsForAccord(accord.ID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}

	if tags[0] != "fresh" {
		t.Errorf("Expected tag 'fresh', got '%s'", tags[0])
	}
}

func TestAccordRepository_RemoveTag(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	// Add tags
	err = repo.AddTag(accord.ID, "fresh")
	if err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}

	err = repo.AddTag(accord.ID, "citrus")
	if err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}

	// Remove one tag
	err = repo.RemoveTag(accord.ID, "fresh")
	if err != nil {
		t.Fatalf("Failed to remove tag: %v", err)
	}

	// Get remaining tags
	tags, err := repo.GetTagsForAccord(accord.ID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}

	if tags[0] != "citrus" {
		t.Errorf("Expected tag 'citrus', got '%s'", tags[0])
	}
}

func TestAccordRepository_SetTags(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)
	accord := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord)
	if err != nil {
		t.Fatalf("Failed to create accord: %v", err)
	}

	// Set tags (replaces any existing)
	newTags := []string{"fresh", "citrus", "summer"}
	err = repo.SetTags(accord.ID, newTags)
	if err != nil {
		t.Fatalf("Failed to set tags: %v", err)
	}

	// Get tags
	tags, err := repo.GetTagsForAccord(accord.ID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
}

func TestAccordRepository_UniqueConstraint(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := NewAccordRepository(tdb.DB)

	// Create first accord
	accord1 := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        25.5,
	}
	err = repo.Create(accord1)
	if err != nil {
		t.Fatalf("Failed to create first accord: %v", err)
	}

	// Try to create duplicate (same name + position for same user)
	accord2 := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "top",
		VolumeMl:        30.0,
	}
	err = repo.Create(accord2)
	if err == nil {
		t.Error("Expected error for duplicate name+position, got nil")
	}

	// Should succeed with different position
	accord3 := &models.Accord{
		UserID:          user.ID,
		Name:            "Citrus Fresh",
		PyramidPosition: "middle",
		VolumeMl:        30.0,
	}
	err = repo.Create(accord3)
	if err != nil {
		t.Errorf("Expected success for different position, got error: %v", err)
	}
}
