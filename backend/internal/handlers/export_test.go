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

func setupExportHandler(t *testing.T) (*ExportHandler, *services.AccordService, *testutil.TestDB, string) {
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
	handler := NewExportHandler(accordService)

	return handler, accordService, tdb, userID
}

func TestExportHandler_Export(t *testing.T) {
	handler, accordService, tdb, userID := setupExportHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create a test accord
	createReq := &models.CreateAccordRequest{
		Name:            "Test Accord",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	}
	_, err := accordService.CreateAccord(userID, createReq)
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err = handler.Export(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var exportData models.ExportData
	err = json.Unmarshal(rec.Body.Bytes(), &exportData)
	require.NoError(t, err)
	assert.Equal(t, "1.0", exportData.Version)
	assert.Len(t, exportData.Accords, 1)
	assert.Equal(t, "Test Accord", exportData.Accords[0].Name)
}

func TestExportHandler_Import(t *testing.T) {
	handler, _, tdb, userID := setupExportHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	importData := models.ImportData{
		Accords: []models.ImportAccord{
			{
				Name:            "Imported Accord 1",
				PyramidPosition: "top",
				VolumeMl:        5.0,
			},
			{
				Name:            "Imported Accord 2",
				PyramidPosition: "middle",
				VolumeMl:        10.0,
			},
		},
	}
	body, _ := json.Marshal(importData)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/import", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.Import(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.ImportResult2
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalRecords)
	assert.Equal(t, 2, result.ImportedRecords, "Errors: %v", result.Errors)
	assert.Equal(t, 0, result.FailedRecords)
}
