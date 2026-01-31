package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/services"
	"github.com/yourusername/scentora-backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func setupInvitationHandler(t *testing.T) (*InvitationHandler, *services.InvitationService, *testutil.TestDB, string) {
	tdb := testutil.SetupTestDB(t)

	// Create test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	var userID string
	err = tdb.DB.QueryRow(`
		INSERT INTO users (email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`, "test@test.com", "testuser", string(hashedPassword)).Scan(&userID)
	require.NoError(t, err)

	// Setup services
	invitationRepo := repository.NewInvitationRepository(tdb.DB)
	invitationService := services.NewInvitationService(invitationRepo)
	handler := NewInvitationHandler(invitationService)

	return handler, invitationService, tdb, userID
}

func TestInvitationHandler_Create(t *testing.T) {
	handler, _, tdb, userID := setupInvitationHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	email := "invitee@test.com"
	reqBody := models.CreateInvitationRequest{
		Email:         &email,
		ExpiresInDays: 7,
	}
	body, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/invitations", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.Create(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "invitation")
}

func TestInvitationHandler_Create_Unauthorized(t *testing.T) {
	handler, _, tdb, _ := setupInvitationHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	email := "invitee@test.com"
	reqBody := models.CreateInvitationRequest{
		Email:         &email,
		ExpiresInDays: 7,
	}
	body, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/invitations", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No userId set - should be unauthorized

	err := handler.Create(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestInvitationHandler_List(t *testing.T) {
	handler, service, tdb, userID := setupInvitationHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create some invitations
	email1 := "test1@test.com"
	email2 := "test2@test.com"
	_, _ = service.Create(userID, &email1, 7)
	_, _ = service.Create(userID, &email2, 7)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/invitations", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.List(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "invitations")
}

func TestInvitationHandler_Revoke(t *testing.T) {
	handler, service, tdb, userID := setupInvitationHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create invitation
	email := "test@test.com"
	invitation, _ := service.Create(userID, &email, 7)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/invitations/"+invitation.Code+"/revoke", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/invitations/:code/revoke")
	c.SetParamNames("code")
	c.SetParamValues(invitation.Code)
	c.Set("userId", userID)

	err := handler.Revoke(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
