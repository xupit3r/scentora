package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRefreshTokenRepository_Create(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user
	user := &models.User{
		Email:        "token_test@example.com",
		Username:     "tokenuser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create a refresh token
	token := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_abc123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = tokenRepo.Create(token)
	require.NoError(t, err)
	assert.NotEmpty(t, token.ID)
	assert.NotZero(t, token.CreatedAt)
}

func TestRefreshTokenRepository_FindByTokenHash(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user
	user := &models.User{
		Email:        "find_token@example.com",
		Username:     "findtokenuser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create a refresh token
	token := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_def456",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token)
	require.NoError(t, err)

	// Find the token
	found, err := tokenRepo.FindByTokenHash("hash_def456")
	require.NoError(t, err)
	assert.Equal(t, token.ID, found.ID)
	assert.Equal(t, user.ID, found.UserID)
	assert.Equal(t, "hash_def456", found.TokenHash)
	assert.False(t, found.Revoked)
}

func TestRefreshTokenRepository_FindByTokenHashNotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Try to find non-existent token
	_, err := tokenRepo.FindByTokenHash("nonexistent_hash")
	assert.Error(t, err)
}

func TestRefreshTokenRepository_FindByTokenHashRevoked(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user
	user := &models.User{
		Email:        "revoked_token@example.com",
		Username:     "revokeduser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create a refresh token
	token := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_ghi789",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token)
	require.NoError(t, err)

	// Revoke the token
	err = tokenRepo.Revoke("hash_ghi789")
	require.NoError(t, err)

	// Try to find the revoked token (should not be found)
	_, err = tokenRepo.FindByTokenHash("hash_ghi789")
	assert.Error(t, err, "Should not find revoked tokens")
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user
	user := &models.User{
		Email:        "revoke_test@example.com",
		Username:     "revokeuser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create a refresh token
	token := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_jkl012",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token)
	require.NoError(t, err)

	// Verify token is not revoked
	found, err := tokenRepo.FindByTokenHash("hash_jkl012")
	require.NoError(t, err)
	assert.False(t, found.Revoked)

	// Revoke the token
	err = tokenRepo.Revoke("hash_jkl012")
	require.NoError(t, err)

	// Verify token is now revoked (not found by FindByTokenHash)
	_, err = tokenRepo.FindByTokenHash("hash_jkl012")
	assert.Error(t, err, "Revoked tokens should not be found")
}

func TestRefreshTokenRepository_RevokeNotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Try to revoke non-existent token
	err := tokenRepo.Revoke("nonexistent_hash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token not found")
}

func TestRefreshTokenRepository_RevokeAllForUser(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user
	user := &models.User{
		Email:        "revokeall@example.com",
		Username:     "revokealluser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create multiple refresh tokens for the user
	token1 := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_token1",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token1)
	require.NoError(t, err)

	token2 := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_token2",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token2)
	require.NoError(t, err)

	token3 := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hash_token3",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token3)
	require.NoError(t, err)

	// Verify all tokens exist
	_, err = tokenRepo.FindByTokenHash("hash_token1")
	require.NoError(t, err)
	_, err = tokenRepo.FindByTokenHash("hash_token2")
	require.NoError(t, err)
	_, err = tokenRepo.FindByTokenHash("hash_token3")
	require.NoError(t, err)

	// Revoke all tokens for the user
	err = tokenRepo.RevokeAllForUser(user.ID)
	require.NoError(t, err)

	// Verify all tokens are revoked
	_, err = tokenRepo.FindByTokenHash("hash_token1")
	assert.Error(t, err, "Token 1 should be revoked")
	_, err = tokenRepo.FindByTokenHash("hash_token2")
	assert.Error(t, err, "Token 2 should be revoked")
	_, err = tokenRepo.FindByTokenHash("hash_token3")
	assert.Error(t, err, "Token 3 should be revoked")
}

func TestRefreshTokenRepository_RevokeAllForUserNoTokens(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create a test user with no tokens
	user := &models.User{
		Email:        "notokens@example.com",
		Username:     "notokensuser",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Revoke all tokens for user (should not error even with no tokens)
	err = tokenRepo.RevokeAllForUser(user.ID)
	assert.NoError(t, err)
}

func TestRefreshTokenRepository_MultipleUsers(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	userRepo := NewUserRepository(tdb.DB)
	tokenRepo := NewRefreshTokenRepository(tdb.DB)

	// Create two test users
	user1 := &models.User{
		Email:        "user1@example.com",
		Username:     "user1",
		PasswordHash: "hashed_password",
	}
	err := userRepo.Create(user1)
	require.NoError(t, err)

	user2 := &models.User{
		Email:        "user2@example.com",
		Username:     "user2",
		PasswordHash: "hashed_password",
	}
	err = userRepo.Create(user2)
	require.NoError(t, err)

	// Create tokens for both users
	token1 := &models.RefreshToken{
		UserID:    user1.ID,
		TokenHash: "hash_user1_token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token1)
	require.NoError(t, err)

	token2 := &models.RefreshToken{
		UserID:    user2.ID,
		TokenHash: "hash_user2_token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = tokenRepo.Create(token2)
	require.NoError(t, err)

	// Revoke all tokens for user1
	err = tokenRepo.RevokeAllForUser(user1.ID)
	require.NoError(t, err)

	// Verify user1's token is revoked
	_, err = tokenRepo.FindByTokenHash("hash_user1_token")
	assert.Error(t, err, "User1's token should be revoked")

	// Verify user2's token is still active
	found, err := tokenRepo.FindByTokenHash("hash_user2_token")
	require.NoError(t, err)
	assert.Equal(t, user2.ID, found.UserID)
	assert.False(t, found.Revoked)
}
