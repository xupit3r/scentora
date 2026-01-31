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

func setupAccordService(t *testing.T) (*AccordService, *testutil.TestDB, string) {
	tdb := testutil.SetupTestDB(t)
	accordRepo := repository.NewAccordRepository(tdb.DB)
	tagRepo := repository.NewPredefinedTagRepository(tdb.DB)
	service := NewAccordService(accordRepo, tagRepo)
	
	// Create a test user
	user := createTestUser(t, tdb, "accord@example.com", "accorduser", "password")
	
	return service, tdb, user.ID
}

func TestAccordService_CreateAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateAccordRequest{
		Name:            "Bergamot Essential Oil",
		PyramidPosition: "top",
		VolumeMl:        15.5,
		Tags:            []string{"citrus", "fresh"},
	}

	accord, err := service.CreateAccord(userID, req)

	require.NoError(t, err)
	require.NotNil(t, accord)
	assert.NotEmpty(t, accord.ID)
	assert.Equal(t, userID, accord.UserID)
	assert.Equal(t, "Bergamot Essential Oil", accord.Name)
	assert.Equal(t, "top", accord.PyramidPosition)
	assert.Equal(t, 15.5, accord.VolumeMl)
}

func TestAccordService_CreateAccord_WithAllFields(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	supplier := "Eden Botanicals"
	purchaseDate := time.Now()
	dilution := 10.0
	notes := "Premium quality, organic"

	req := &models.CreateAccordRequest{
		Name:               "Rose Otto",
		PyramidPosition:    "middle",
		VolumeMl:           5.0,
		Supplier:           &supplier,
		PurchaseDate:       &purchaseDate,
		DilutionPercentage: &dilution,
		Notes:              &notes,
		Tags:               []string{"floral", "elegant"},
	}

	accord, err := service.CreateAccord(userID, req)

	require.NoError(t, err)
	assert.NotNil(t, accord.Supplier)
	assert.Equal(t, supplier, *accord.Supplier)
	assert.NotNil(t, accord.PurchaseDate)
	assert.NotNil(t, accord.DilutionPercentage)
	assert.Equal(t, dilution, *accord.DilutionPercentage)
	assert.NotNil(t, accord.Notes)
	assert.Equal(t, notes, *accord.Notes)
}

func TestAccordService_CreateAccord_MissingName(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateAccordRequest{
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}

	_, err := service.CreateAccord(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestAccordService_CreateAccord_MissingPyramidPosition(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateAccordRequest{
		Name:     "Test Accord",
		VolumeMl: 10.0,
	}

	_, err := service.CreateAccord(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pyramid position is required")
}

func TestAccordService_CreateAccord_InvalidPyramidPosition(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateAccordRequest{
		Name:            "Test Accord",
		PyramidPosition: "invalid",
		VolumeMl:        10.0,
	}

	_, err := service.CreateAccord(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be one of: top, middle, base")
}

func TestAccordService_CreateAccord_InvalidVolume(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	req := &models.CreateAccordRequest{
		Name:            "Test Accord",
		PyramidPosition: "top",
		VolumeMl:        0,
	}

	_, err := service.CreateAccord(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "volume must be greater than 0")
}

func TestAccordService_CreateAccord_InvalidDilution(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	invalidDilution := 150.0
	req := &models.CreateAccordRequest{
		Name:               "Test Accord",
		PyramidPosition:    "top",
		VolumeMl:           10.0,
		DilutionPercentage: &invalidDilution,
	}

	_, err := service.CreateAccord(userID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dilution percentage must be between 0 and 100")
}

func TestAccordService_GetAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	req := &models.CreateAccordRequest{
		Name:            "Lavender",
		PyramidPosition: "middle",
		VolumeMl:        20.0,
		Tags:            []string{"floral", "calming"},
	}
	created, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Get the accord
	fetched, err := service.GetAccord(created.ID, userID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Lavender", fetched.Name)
	assert.Len(t, fetched.Tags, 2)
	assert.Contains(t, fetched.Tags, "floral")
	assert.Contains(t, fetched.Tags, "calming")
}

func TestAccordService_GetAccord_NotFound(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	_, err := service.GetAccord("00000000-0000-0000-0000-000000000000", userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAccordService_GetAccord_WrongUser(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create another user
	otherUser := createTestUser(t, tdb, "other@example.com", "other", "password")

	// User1 creates an accord
	req := &models.CreateAccordRequest{
		Name:            "User1 Accord",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// User2 tries to get user1's accord
	_, err = service.GetAccord(accord.ID, otherUser.ID)
	assert.Error(t, err)
}

func TestAccordService_ListAccords(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create multiple accords
	accords := []struct {
		name     string
		position string
		volume   float64
	}{
		{"Bergamot", "top", 15.0},
		{"Rose", "middle", 10.0},
		{"Sandalwood", "base", 20.0},
	}

	for _, a := range accords {
		req := &models.CreateAccordRequest{
			Name:            a.name,
			PyramidPosition: a.position,
			VolumeMl:        a.volume,
		}
		_, err := service.CreateAccord(userID, req)
		require.NoError(t, err)
	}

	// List accords
	list, err := service.ListAccords(userID)

	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestAccordService_ListAccords_Empty(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	list, err := service.ListAccords(userID)

	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestAccordService_ListAccords_UserIsolation(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create another user
	otherUser := createTestUser(t, tdb, "other@example.com", "other", "password")

	// User1 creates 2 accords
	req1 := &models.CreateAccordRequest{
		Name:            "User1 Accord1",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	_, err := service.CreateAccord(userID, req1)
	require.NoError(t, err)

	req2 := &models.CreateAccordRequest{
		Name:            "User1 Accord2",
		PyramidPosition: "middle",
		VolumeMl:        15.0,
	}
	_, err = service.CreateAccord(userID, req2)
	require.NoError(t, err)

	// User2 creates 1 accord
	req3 := &models.CreateAccordRequest{
		Name:            "User2 Accord",
		PyramidPosition: "base",
		VolumeMl:        20.0,
	}
	_, err = service.CreateAccord(otherUser.ID, req3)
	require.NoError(t, err)

	// List user1's accords
	list1, err := service.ListAccords(userID)
	require.NoError(t, err)
	assert.Len(t, list1, 2)

	// List user2's accords
	list2, err := service.ListAccords(otherUser.ID)
	require.NoError(t, err)
	assert.Len(t, list2, 1)
}

func TestAccordService_UpdateAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	req := &models.CreateAccordRequest{
		Name:            "Original Name",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Update the accord
	newName := "Updated Name"
	newVolume := 15.5
	updateReq := &models.UpdateAccordRequest{
		Name:     &newName,
		VolumeMl: &newVolume,
	}

	updated, err := service.UpdateAccord(accord.ID, userID, updateReq)

	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newVolume, updated.VolumeMl)
	assert.Equal(t, "top", updated.PyramidPosition, "Unchanged fields should remain")
}

func TestAccordService_UpdateAccord_Tags(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create accord with initial tags
	req := &models.CreateAccordRequest{
		Name:            "Test Accord",
		PyramidPosition: "top",
		VolumeMl:        10.0,
		Tags:            []string{"tag1", "tag2"},
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Update tags
	newTags := []string{"tag3", "tag4", "tag5"}
	updateReq := &models.UpdateAccordRequest{
		Tags: &newTags,
	}

	updated, err := service.UpdateAccord(accord.ID, userID, updateReq)

	require.NoError(t, err)
	assert.Len(t, updated.Tags, 3)
	assert.Contains(t, updated.Tags, "tag3")
	assert.Contains(t, updated.Tags, "tag4")
	assert.Contains(t, updated.Tags, "tag5")
	assert.NotContains(t, updated.Tags, "tag1")
}

func TestAccordService_UpdateAccord_InvalidVolume(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	req := &models.CreateAccordRequest{
		Name:            "Test",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Try to update with invalid volume
	invalidVolume := -5.0
	updateReq := &models.UpdateAccordRequest{
		VolumeMl: &invalidVolume,
	}

	_, err = service.UpdateAccord(accord.ID, userID, updateReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "volume must be greater than 0")
}

func TestAccordService_DeleteAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	req := &models.CreateAccordRequest{
		Name:            "To Delete",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Delete the accord
	err = service.DeleteAccord(accord.ID, userID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = service.GetAccord(accord.ID, userID)
	assert.Error(t, err)
}

func TestAccordService_DeleteAccord_NotFound(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	err := service.DeleteAccord("00000000-0000-0000-0000-000000000000", userID)
	assert.Error(t, err)
}

func TestAccordService_AddTagToAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create accord without tags
	req := &models.CreateAccordRequest{
		Name:            "Test",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Add a tag
	err = service.AddTagToAccord(accord.ID, userID, "newtag")
	require.NoError(t, err)

	// Verify tag was added
	fetched, err := service.GetAccord(accord.ID, userID)
	require.NoError(t, err)
	assert.Contains(t, fetched.Tags, "newtag")
}

func TestAccordService_RemoveTagFromAccord(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create accord with tags
	req := &models.CreateAccordRequest{
		Name:            "Test",
		PyramidPosition: "top",
		VolumeMl:        10.0,
		Tags:            []string{"tag1", "tag2", "tag3"},
	}
	accord, err := service.CreateAccord(userID, req)
	require.NoError(t, err)

	// Remove a tag
	err = service.RemoveTagFromAccord(accord.ID, userID, "tag2")
	require.NoError(t, err)

	// Verify tag was removed
	fetched, err := service.GetAccord(accord.ID, userID)
	require.NoError(t, err)
	assert.Len(t, fetched.Tags, 2)
	assert.Contains(t, fetched.Tags, "tag1")
	assert.Contains(t, fetched.Tags, "tag3")
	assert.NotContains(t, fetched.Tags, "tag2")
}

func TestAccordService_PyramidPositions(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	positions := []string{"top", "middle", "base"}

	for _, pos := range positions {
		req := &models.CreateAccordRequest{
			Name:            pos + " Note",
			PyramidPosition: pos,
			VolumeMl:        10.0,
		}
		accord, err := service.CreateAccord(userID, req)
		require.NoError(t, err, "Should create accord with position: %s", pos)
		assert.Equal(t, pos, accord.PyramidPosition)
	}
}

func TestAccordService_DuplicateNameSamePosition(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create first accord
	req1 := &models.CreateAccordRequest{
		Name:            "Duplicate",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	_, err := service.CreateAccord(userID, req1)
	require.NoError(t, err)

	// Try to create duplicate (same name, same position)
	req2 := &models.CreateAccordRequest{
		Name:            "Duplicate",
		PyramidPosition: "top",
		VolumeMl:        15.0,
	}
	_, err = service.CreateAccord(userID, req2)
	assert.Error(t, err, "Should not allow duplicate name in same position")
}

func TestAccordService_SameNameDifferentPosition(t *testing.T) {
	service, tdb, userID := setupAccordService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create accord in top position
	req1 := &models.CreateAccordRequest{
		Name:            "Rose",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	accord1, err := service.CreateAccord(userID, req1)
	require.NoError(t, err)

	// Create accord with same name in different position (should work)
	req2 := &models.CreateAccordRequest{
		Name:            "Rose",
		PyramidPosition: "middle",
		VolumeMl:        15.0,
	}
	accord2, err := service.CreateAccord(userID, req2)
	require.NoError(t, err)

	// Verify both exist and are different
	assert.NotEqual(t, accord1.ID, accord2.ID)
	assert.Equal(t, "top", accord1.PyramidPosition)
	assert.Equal(t, "middle", accord2.PyramidPosition)
}
