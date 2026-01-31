package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func setupInvitationService(t *testing.T) (*InvitationService, *testutil.TestDB) {
	tdb := testutil.SetupTestDB(t)
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	service := NewInvitationService(invitationRepo)
	return service, tdb
}

func TestInvitationService_Create(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an invitation
	invitation, err := service.Create(creator.ID, nil, 7)

	require.NoError(t, err)
	require.NotNil(t, invitation)
	assert.NotEmpty(t, invitation.ID)
	assert.NotEmpty(t, invitation.Code)
	assert.Equal(t, creator.ID, invitation.CreatedBy)
	assert.False(t, invitation.Used)
	assert.Nil(t, invitation.Email)
	assert.True(t, invitation.ExpiresAt.After(time.Now()))
}

func TestInvitationService_Create_WithEmail(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an email-specific invitation
	email := "invited@example.com"
	invitation, err := service.Create(creator.ID, &email, 14)

	require.NoError(t, err)
	require.NotNil(t, invitation)
	assert.NotEmpty(t, invitation.Code)
	assert.Equal(t, creator.ID, invitation.CreatedBy)
	assert.NotNil(t, invitation.Email)
	assert.Equal(t, email, *invitation.Email)
	assert.False(t, invitation.Used)
}

func TestInvitationService_Create_UniqueCode(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create multiple invitations
	invitation1, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	invitation2, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	invitation3, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	// Verify all codes are unique
	assert.NotEqual(t, invitation1.Code, invitation2.Code)
	assert.NotEqual(t, invitation1.Code, invitation3.Code)
	assert.NotEqual(t, invitation2.Code, invitation3.Code)
}

func TestInvitationService_List(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create two users
	creator1 := createTestUser(t, tdb, "creator1@example.com", "creator1", "password")
	creator2 := createTestUser(t, tdb, "creator2@example.com", "creator2", "password")

	// Creator1 creates 3 invitations
	_, err := service.Create(creator1.ID, nil, 7)
	require.NoError(t, err)
	_, err = service.Create(creator1.ID, nil, 7)
	require.NoError(t, err)
	_, err = service.Create(creator1.ID, nil, 7)
	require.NoError(t, err)

	// Creator2 creates 2 invitations
	_, err = service.Create(creator2.ID, nil, 7)
	require.NoError(t, err)
	_, err = service.Create(creator2.ID, nil, 7)
	require.NoError(t, err)

	// List creator1's invitations
	invitations1, err := service.List(creator1.ID)
	require.NoError(t, err)
	assert.Len(t, invitations1, 3)

	// List creator2's invitations
	invitations2, err := service.List(creator2.ID)
	require.NoError(t, err)
	assert.Len(t, invitations2, 2)

	// Verify isolation
	for _, inv := range invitations1 {
		assert.Equal(t, creator1.ID, inv.CreatedBy)
	}
	for _, inv := range invitations2 {
		assert.Equal(t, creator2.ID, inv.CreatedBy)
	}
}

func TestInvitationService_List_Empty(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a user who hasn't created any invitations
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// List invitations
	invitations, err := service.List(creator.ID)
	require.NoError(t, err)
	assert.Empty(t, invitations)
}

func TestInvitationService_Revoke(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an invitation
	invitation, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	// Verify invitation exists
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	found, err := invitationRepo.FindByCode(invitation.Code)
	require.NoError(t, err)
	assert.False(t, found.Used)

	// Revoke the invitation
	err = service.Revoke(invitation.Code, creator.ID)
	require.NoError(t, err)

	// Verify invitation is marked as used (revoked)
	found, err = invitationRepo.FindByCode(invitation.Code)
	require.NoError(t, err)
	assert.True(t, found.Used)
}

func TestInvitationService_Revoke_WrongCreator(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create two users
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	otherUser := createTestUser(t, tdb, "other@example.com", "other", "password")

	// Creator creates an invitation
	invitation, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	// OtherUser tries to revoke creator's invitation
	err = service.Revoke(invitation.Code, otherUser.ID)
	assert.Error(t, err, "Should not allow revoking another user's invitation")

	// Verify invitation is still active
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	found, err := invitationRepo.FindByCode(invitation.Code)
	require.NoError(t, err)
	assert.False(t, found.Used)
}

func TestInvitationService_ValidateAndUse(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an invitation
	invitation, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	// Create a new user
	newUser := createTestUser(t, tdb, "newuser@example.com", "newuser", "password")

	// Validate and use the invitation
	err = service.ValidateAndUse(invitation.Code, newUser.Email, newUser.ID)
	require.NoError(t, err)

	// Verify invitation is marked as used
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	found, err := invitationRepo.FindByCode(invitation.Code)
	require.NoError(t, err)
	assert.True(t, found.Used)
	assert.NotNil(t, found.UsedBy)
	assert.Equal(t, newUser.ID, *found.UsedBy)
}

func TestInvitationService_ValidateAndUse_InvalidCode(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Try to use invalid invitation code
	err := service.ValidateAndUse("INVALID-CODE", "test@example.com", "user-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid invitation")
}

func TestInvitationService_ValidateAndUse_AlreadyUsed(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an invitation
	invitation, err := service.Create(creator.ID, nil, 7)
	require.NoError(t, err)

	// Create first user
	user1 := createTestUser(t, tdb, "user1@example.com", "user1", "password")

	// First use
	err = service.ValidateAndUse(invitation.Code, user1.Email, user1.ID)
	require.NoError(t, err)

	// Create second user
	user2 := createTestUser(t, tdb, "user2@example.com", "user2", "password")

	// Try to use again
	err = service.ValidateAndUse(invitation.Code, user2.Email, user2.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already been used")
}

func TestInvitationService_ValidateAndUse_Expired(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an expired invitation manually
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	expiredInvitation := &models.Invitation{
		Code:      "EXPIRED-CODE",
		CreatedBy: creator.ID,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired yesterday
		Used:      false,
	}
	err := invitationRepo.Create(expiredInvitation)
	require.NoError(t, err)

	// Try to use expired invitation
	err = service.ValidateAndUse(expiredInvitation.Code, "test@example.com", "user-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestInvitationService_ValidateAndUse_EmailSpecific(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an email-specific invitation
	specificEmail := "specific@example.com"
	invitation, err := service.Create(creator.ID, &specificEmail, 7)
	require.NoError(t, err)

	// Create a user with the correct email
	user := createTestUser(t, tdb, specificEmail, "user", "password")

	// Use invitation with correct email
	err = service.ValidateAndUse(invitation.Code, specificEmail, user.ID)
	assert.NoError(t, err)
}

func TestInvitationService_ValidateAndUse_WrongEmail(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Create an email-specific invitation
	specificEmail := "specific@example.com"
	invitation, err := service.Create(creator.ID, &specificEmail, 7)
	require.NoError(t, err)

	// Create a user with a different email
	user := createTestUser(t, tdb, "wrong@example.com", "user", "password")

	// Try to use invitation with wrong email
	err = service.ValidateAndUse(invitation.Code, user.Email, user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different email")
}

func TestInvitationService_ExpirationDates(t *testing.T) {
	service, tdb := setupInvitationService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")

	// Test different expiration periods
	testCases := []struct {
		days         int
		expectedDays int
	}{
		{1, 1},
		{7, 7},
		{14, 14},
		{30, 30},
	}

	for _, tc := range testCases {
		invitation, err := service.Create(creator.ID, nil, tc.days)
		require.NoError(t, err)

		// Calculate expected expiration date
		expectedExpiration := time.Now().AddDate(0, 0, tc.expectedDays)

		// Allow for small timing differences (up to 1 second)
		timeDiff := invitation.ExpiresAt.Sub(expectedExpiration)
		assert.True(t, timeDiff < time.Second && timeDiff > -time.Second,
			"Expiration should be approximately %d days from now", tc.expectedDays)
	}
}

// Helper function to create a test user (reused from auth_service_test.go)
func createTestUserForInvitation(t *testing.T, tdb *testutil.TestDB, email, username, password string) *models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	userRepo := repository.NewUserRepository(tdb.DB)
	err = userRepo.Create(user)
	require.NoError(t, err)

	return user
}
