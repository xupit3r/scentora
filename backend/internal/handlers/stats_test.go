package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/services"
	"github.com/yourusername/scentora-backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func setupStatsHandler(t *testing.T) (*StatsHandler, *services.AccordService, *testutil.TestDB, string) {
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

	// Setup repositories and services
	accordRepo := repository.NewAccordRepository(tdb.DB)
	tagRepo := repository.NewPredefinedTagRepository(tdb.DB)
	accordService := services.NewAccordService(accordRepo, tagRepo)
	handler := NewStatsHandler(accordService)

	return handler, accordService, tdb, userID
}

func TestStatsHandler_GetStats(t *testing.T) {
	handler, _, tdb, userID := setupStatsHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.GetStats(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
