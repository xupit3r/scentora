package handlers

import (
"encoding/json"
"net/http"
"net/http/httptest"
"strings"
"testing"

"github.com/labstack/echo/v4"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
"github.com/yourusername/scentora-backend/internal/config"
"github.com/yourusername/scentora-backend/internal/models"
"github.com/yourusername/scentora-backend/internal/repository"
"github.com/yourusername/scentora-backend/internal/services"
"github.com/yourusername/scentora-backend/internal/testutil"
"golang.org/x/crypto/bcrypt"
)

func setupAuthHandlerAndService(t *testing.T) (*AuthHandler, *services.AuthService, *testutil.TestDB) {
tdb := testutil.SetupTestDB(t)

cfg := &config.Config{
JWTSecret:           "test-secret-key",
JWTAccessExpiresIn:  "15m",
JWTRefreshExpiresIn: "7d",
}

userRepo := repository.NewUserRepository(tdb.DB)
tokenRepo := repository.NewRefreshTokenRepository(tdb.DB)
invitationRepo := repository.NewInvitationRepository(tdb.DB)

authService := services.NewAuthService(userRepo, tokenRepo, invitationRepo, cfg)
handler := NewAuthHandler(authService)

return handler, authService, tdb
}

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

func TestAuthHandler_Register(t *testing.T) {
handler, _, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create a creator user and invitation
creator := createTestUser(t, tdb, "creator@example.com", "creator", "password")
invitationService := services.NewInvitationService(repository.NewInvitationRepository(tdb.DB))
invitation, err := invitationService.Create(creator.ID, nil, 7)
require.NoError(t, err)

// Setup Echo
e := echo.New()
reqBody := `{
"email": "newuser@example.com",
"username": "newuser",
"password": "password123",
"invitationCode": "` + invitation.Code + `"
}`
req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(reqBody))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

// Execute
err = handler.Register(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusCreated, rec.Code)

var response models.AuthResponse
err = json.Unmarshal(rec.Body.Bytes(), &response)
require.NoError(t, err)
assert.NotEmpty(t, response.AccessToken)
assert.NotEmpty(t, response.RefreshToken)
assert.Equal(t, "newuser@example.com", response.User.Email)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
handler, _, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

e := echo.New()
req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader("invalid json"))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

err := handler.Register(c)

require.NoError(t, err)
assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Login(t *testing.T) {
handler, _, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create test user
password := "password123"
user := createTestUser(t, tdb, "login", "loginuser", password)

// Setup Echo
e := echo.New()
reqBody := `{"email": "` + user.Email + `", "password": "` + password + `"}`
req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(reqBody))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

// Execute
err := handler.Login(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusOK, rec.Code)

var response models.AuthResponse
err = json.Unmarshal(rec.Body.Bytes(), &response)
require.NoError(t, err)
assert.NotEmpty(t, response.AccessToken)
assert.NotEmpty(t, response.RefreshToken)
assert.Equal(t, user.Email, response.User.Email)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
handler, _, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create test user
createTestUser(t, tdb, "test@example.com", "testuser", "correctpassword")

// Setup Echo with wrong password
e := echo.New()
reqBody := `{"email": "test@example.com", "password": "wrongpassword"}`
req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(reqBody))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

// Execute
err := handler.Login(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_Refresh(t *testing.T) {
handler, authService, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create and login user
password := "password123"
user := createTestUser(t, tdb, "refresh", "refreshuser", password)

loginResponse, err := authService.Login(user.Email, password)
require.NoError(t, err)

// Setup Echo
e := echo.New()
reqBody := `{"refreshToken": "` + loginResponse.RefreshToken + `"}`
req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(reqBody))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

// Execute
err = handler.Refresh(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusOK, rec.Code)

var response models.AuthResponse
err = json.Unmarshal(rec.Body.Bytes(), &response)
require.NoError(t, err)
assert.NotEmpty(t, response.AccessToken)
assert.NotEmpty(t, response.RefreshToken)
}

func TestAuthHandler_Logout(t *testing.T) {
handler, authService, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create and login user
password := "password123"
user := createTestUser(t, tdb, "logout", "logoutuser", password)

loginResponse, err := authService.Login(user.Email, password)
require.NoError(t, err)

// Setup Echo
e := echo.New()
reqBody := `{"refreshToken": "` + loginResponse.RefreshToken + `"}`
req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", strings.NewReader(reqBody))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)

// Execute
err = handler.Logout(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusNoContent, rec.Code)

// Verify token is revoked
_, err = authService.Refresh(loginResponse.RefreshToken)
assert.Error(t, err)
}

func TestAuthHandler_Me(t *testing.T) {
handler, _, tdb := setupAuthHandlerAndService(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create test user
user := createTestUser(t, tdb, "me@example.com", "meuser", "password")

// Setup Echo with JWT middleware context
e := echo.New()
req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)
c.Set("userId", user.ID)

// Execute
err := handler.Me(c)

// Assert
require.NoError(t, err)
assert.Equal(t, http.StatusOK, rec.Code)

var responseUser models.User
err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
require.NoError(t, err)
assert.Equal(t, user.Email, responseUser.Email)
assert.Empty(t, responseUser.PasswordHash)
}
