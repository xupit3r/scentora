package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/config"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to create auth service with test database
func setupAuthService(t *testing.T) (*AuthService, *testutil.TestDB) {
	tdb := testutil.SetupTestDB(t)
	
	cfg := &config.Config{
		JWTSecret:           "test-secret-key-for-testing-only",
		JWTAccessExpiresIn:  "15m",
		JWTRefreshExpiresIn: "7d",
	}
	
	userRepo := repository.NewUserRepository(tdb.DB)
	tokenRepo := repository.NewRefreshTokenRepository(tdb.DB)
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	
	authService := NewAuthService(userRepo, tokenRepo, invitationRepo, cfg)
	
	return authService, tdb
}

// Helper function to create a test user directly in database
func createTestUser(t *testing.T, tdb *testutil.TestDB, emailPrefix, usernamePrefix, password string) *models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	
	user := &models.User{
		Email:        testutil.UniqueEmail(emailPrefix),
		Username:     usernamePrefix + "_" + testutil.UniqueEmail(""),
		PasswordHash: string(hashedPassword),
	}
	
	userRepo := repository.NewUserRepository(tdb.DB)
	err = userRepo.Create(user)
	require.NoError(t, err)
	
	return user
}

// Helper function to create a valid invitation
func createTestInvitation(t *testing.T, tdb *testutil.TestDB, createdByUserID string, email *string, expiresInDays int) *models.Invitation {
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	
	invitation := &models.Invitation{
		Code:      "TEST-INVITE-" + time.Now().Format("20060102150405"),
		Email:     email,
		CreatedBy: createdByUserID,
		ExpiresAt: time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour),
	}
	
	err := invitationRepo.Create(invitation)
	require.NoError(t, err)
	
	return invitation
}

func TestAuthService_Register(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create a test user to create invitation
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	
	// Create a valid invitation
	invitation := createTestInvitation(t, tdb, creator.ID, nil, 7)
	
	// Register a new user
	authResponse, err := authService.Register(
		"newuser@example.com",
		"newuser",
		"password123",
		invitation.Code,
	)
	
	require.NoError(t, err)
	require.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)
	assert.NotNil(t, authResponse.User)
	assert.Equal(t, "newuser@example.com", authResponse.User.Email)
	assert.Equal(t, "newuser", authResponse.User.Username)
	assert.Empty(t, authResponse.User.PasswordHash, "Password hash should not be in response")
}

func TestAuthService_Register_InvalidInvitation(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Try to register with invalid invitation code
	_, err := authService.Register(
		"test@example.com",
		"testuser",
		"password123",
		"INVALID-CODE",
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid invitation")
}

func TestAuthService_Register_ExpiredInvitation(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	
	// Create expired invitation (expires in the past)
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	expiredInvitation := &models.Invitation{
		Code:      "EXPIRED-INVITE",
		CreatedBy: creator.ID,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired yesterday
	}
	err := invitationRepo.Create(expiredInvitation)
	require.NoError(t, err)
	
	// Try to register with expired invitation
	_, err = authService.Register(
		"test@example.com",
		"testuser",
		"password123",
		expiredInvitation.Code,
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestAuthService_Register_UsedInvitation(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	
	// Create invitation
	invitation := createTestInvitation(t, tdb, creator.ID, nil, 7)
	
	// Register first user (uses the invitation)
	_, err := authService.Register(
		"user1@example.com",
		"user1",
		"password123",
		invitation.Code,
	)
	require.NoError(t, err)
	
	// Try to register second user with same invitation
	_, err = authService.Register(
		"user2@example.com",
		"user2",
		"password123",
		invitation.Code,
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already been used")
}

func TestAuthService_Register_EmailSpecificInvitation(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	
	// Create email-specific invitation
	specificEmail := "specific@example.com"
	invitation := createTestInvitation(t, tdb, creator.ID, &specificEmail, 7)
	
	// Try to register with correct email
	authResponse, err := authService.Register(
		specificEmail,
		"user1",
		"password123",
		invitation.Code,
	)
	require.NoError(t, err)
	assert.Equal(t, specificEmail, authResponse.User.Email)
}

func TestAuthService_Register_WrongEmailForEmailSpecificInvitation(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create creator user
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	
	// Create email-specific invitation
	specificEmail := "specific@example.com"
	invitation := createTestInvitation(t, tdb, creator.ID, &specificEmail, 7)
	
	// Try to register with wrong email
	_, err := authService.Register(
		"wrong@example.com",
		"user1",
		"password123",
		invitation.Code,
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different email")
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create existing user
	existingUser := createTestUser(t, tdb, "existing@example.com", "existing", "password")
	
	// Create creator for invitation
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	invitation := createTestInvitation(t, tdb, creator.ID, nil, 7)
	
	// Try to register with duplicate email
	_, err := authService.Register(
		existingUser.Email, // Duplicate email
		"newuser",
		"password123",
		invitation.Code,
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create existing user
	existingUser := createTestUser(t, tdb, "existing@example.com", "existinguser", "password")
	
	// Create creator for invitation
	creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
	invitation := createTestInvitation(t, tdb, creator.ID, nil, 7)
	
	// Try to register with duplicate username
	_, err := authService.Register(
		"newemail@example.com",
		existingUser.Username, // Duplicate username
		"password123",
		invitation.Code,
	)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
}

func TestAuthService_Login(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create test user
	password := "password123"
	user := createTestUser(t, tdb, "login", "loginuser", password)
	
	// Login with the user's actual email
	authResponse, err := authService.Login(user.Email, password)
	
	require.NoError(t, err)
	require.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)
	assert.NotNil(t, authResponse.User)
	assert.Equal(t, user.Email, authResponse.User.Email)
	assert.Empty(t, authResponse.User.PasswordHash)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Try to login with non-existent email
	_, err := authService.Login("nonexistent@example.com", "password")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create test user
	user := createTestUser(t, tdb, "test", "testuser", "correctpassword")
	
	// Try to login with wrong password
	_, err := authService.Login(user.Email, "wrongpassword")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Refresh(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create and login user
	password := "password123"
	user := createTestUser(t, tdb, "refresh", "refreshuser", password)
	
	loginResponse, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	oldAccessToken := loginResponse.AccessToken
	oldRefreshToken := loginResponse.RefreshToken
	
	// Wait a moment to ensure different timestamp
	time.Sleep(time.Second)
	
	// Refresh tokens
	refreshResponse, err := authService.Refresh(oldRefreshToken)
	
	require.NoError(t, err)
	require.NotNil(t, refreshResponse)
	assert.NotEmpty(t, refreshResponse.AccessToken)
	assert.NotEmpty(t, refreshResponse.RefreshToken)
	assert.NotEqual(t, oldAccessToken, refreshResponse.AccessToken, "Access token should be new")
	assert.NotEqual(t, oldRefreshToken, refreshResponse.RefreshToken, "Refresh token should be new")
	assert.Equal(t, user.Email, refreshResponse.User.Email)
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Try to refresh with invalid token
	_, err := authService.Refresh("invalid-refresh-token")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestAuthService_Refresh_RevokedToken(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create and login user
	password := "password123"
	user := createTestUser(t, tdb, "test", "testuser", password)
	
	loginResponse, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	// Logout (revoke refresh token)
	err = authService.Logout(loginResponse.RefreshToken)
	require.NoError(t, err)
	
	// Try to use revoked refresh token
	_, err = authService.Refresh(loginResponse.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestAuthService_Logout(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create and login user
	password := "password123"
	user := createTestUser(t, tdb, "logout", "logoutuser", password)
	
	loginResponse, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	// Logout
	err = authService.Logout(loginResponse.RefreshToken)
	assert.NoError(t, err)
	
	// Verify token is revoked
	_, err = authService.Refresh(loginResponse.RefreshToken)
	assert.Error(t, err)
}

func TestAuthService_LogoutAll(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create and login user multiple times
	password := "password123"
	user := createTestUser(t, tdb, "multidevice", "multiuser", password)
	
	// Login from 3 different "devices"
	login1, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	login2, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	login3, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	// Verify all tokens work
	_, err = authService.Refresh(login1.RefreshToken)
	assert.NoError(t, err)
	_, err = authService.Refresh(login2.RefreshToken)
	assert.NoError(t, err)
	_, err = authService.Refresh(login3.RefreshToken)
	assert.NoError(t, err)
	
	// Logout all
	err = authService.LogoutAll(user.ID)
	assert.NoError(t, err)
	
	// Verify all tokens are revoked
	_, err = authService.Refresh(login1.RefreshToken)
	assert.Error(t, err)
	_, err = authService.Refresh(login2.RefreshToken)
	assert.Error(t, err)
	_, err = authService.Refresh(login3.RefreshToken)
	assert.Error(t, err)
}

func TestAuthService_GetUserByID(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create test user
	user := createTestUser(t, tdb, "getuser", "getuser", "password")
	
	// Get user by ID
	fetchedUser, err := authService.GetUserByID(user.ID)
	
	require.NoError(t, err)
	require.NotNil(t, fetchedUser)
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Username, fetchedUser.Username)
	assert.Empty(t, fetchedUser.PasswordHash, "Password hash should not be returned")
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Try to get non-existent user
	_, err := authService.GetUserByID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

func TestAuthService_TokenGeneration(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create and login user
	password := "password123"
	user := createTestUser(t, tdb, "tokentest", "tokenuser", password)
	
	// Login to generate tokens
	response, err := authService.Login(user.Email, password)
	require.NoError(t, err)
	
	// Verify access token is a valid JWT
	assert.NotEmpty(t, response.AccessToken)
	assert.Contains(t, response.AccessToken, ".")
	
	// Verify refresh token is generated
	assert.NotEmpty(t, response.RefreshToken)
	
	// Verify tokens are different
	assert.NotEqual(t, response.AccessToken, response.RefreshToken)
}

func TestAuthService_PasswordHashing(t *testing.T) {
	authService, tdb := setupAuthService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)
	
	// Create creator and invitation
	creator := createTestUser(t, tdb, "creator", "creator", "password")
	invitation := createTestInvitation(t, tdb, creator.ID, nil, 7)
	
	// Register user
	password := "mysecretpassword"
	_, err := authService.Register(
		"hashtest@example.com",
		"hashuser",
		password,
		invitation.Code,
	)
	require.NoError(t, err)
	
	// Verify password is hashed in database
	userRepo := repository.NewUserRepository(tdb.DB)
	user, err := userRepo.FindByEmail("hashtest@example.com")
	require.NoError(t, err)
	
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, password, user.PasswordHash, "Password should be hashed")
	assert.Contains(t, user.PasswordHash, "$2a$", "Should be bcrypt hash")
	
	// Verify login works with correct password
	_, err = authService.Login("hashtest@example.com", password)
	assert.NoError(t, err)
}
