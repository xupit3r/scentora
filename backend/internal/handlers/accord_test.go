package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/services"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func setupAccordHandler(t *testing.T) (*AccordHandler, *services.AccordService, *testutil.TestDB, string) {
	tdb := testutil.SetupTestDB(t)

	accordRepo := repository.NewAccordRepository(tdb.DB)
	tagRepo := repository.NewPredefinedTagRepository(tdb.DB)
	accordService := services.NewAccordService(accordRepo, tagRepo)
	handler := NewAccordHandler(accordService)

	// Create a test user
	user := createTestUser(t, tdb, "accord@example.com", "accorduser", "password")

	return handler, accordService, tdb, user.ID
}

func TestAccordHandler_Create(t *testing.T) {
	handler, _, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	reqBody := `{
		"name": "Bergamot Essential Oil",
		"pyramidPosition": "top",
		"volumeMl": 15.5,
		"tags": ["citrus", "fresh"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/accords", strings.NewReader(reqBody))
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
	assert.Contains(t, response, "accord")
}

func TestAccordHandler_Create_Unauthorized(t *testing.T) {
	handler, _, tdb, _ := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	reqBody := `{"name": "Test", "pyramidPosition": "top", "volumeMl": 10}`
	req := httptest.NewRequest(http.MethodPost, "/api/accords", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No userId set

	err := handler.Create(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAccordHandler_Create_InvalidJSON(t *testing.T) {
	handler, _, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/accords", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.Create(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAccordHandler_Get(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord first
	accord, err := service.CreateAccord(userID, &models.CreateAccordRequest{
		Name:            "Test Accord",
		PyramidPosition: "middle",
		VolumeMl:        20.0,
	})
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/accords/"+accord.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id")
	c.SetParamNames("id")
	c.SetParamValues(accord.ID)
	c.Set("userId", userID)

	err = handler.Get(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "accord")
}

func TestAccordHandler_Get_NotFound(t *testing.T) {
	handler, _, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/accords/nonexistent-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id")
	c.SetParamNames("id")
	c.SetParamValues("00000000-0000-0000-0000-000000000000")
	c.Set("userId", userID)

	err := handler.Get(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAccordHandler_List(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create multiple accords
	for i := 1; i <= 3; i++ {
		_, err := service.CreateAccord(userID, &models.CreateAccordRequest{
			Name:            fmt.Sprintf("Accord %d", i),
			PyramidPosition: "top",
			VolumeMl:        float64(i * 10),
		})
		require.NoError(t, err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/accords", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", userID)

	err := handler.List(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "accords")
}

func TestAccordHandler_Update(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	accord, err := service.CreateAccord(userID, &models.CreateAccordRequest{
		Name:            "Original Name",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	})
	require.NoError(t, err)

	e := echo.New()
	reqBody := `{"name": "Updated Name", "volumeMl": 15.5}`
	req := httptest.NewRequest(http.MethodPut, "/api/accords/"+accord.ID, strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id")
	c.SetParamNames("id")
	c.SetParamValues(accord.ID)
	c.Set("userId", userID)

	err = handler.Update(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "accord")
}

func TestAccordHandler_Delete(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	accord, err := service.CreateAccord(userID, &models.CreateAccordRequest{
		Name:            "To Delete",
		PyramidPosition: "base",
		VolumeMl:        5.0,
	})
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/accords/"+accord.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id")
	c.SetParamNames("id")
	c.SetParamValues(accord.ID)
	c.Set("userId", userID)

	err = handler.Delete(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify it's deleted
	_, err = service.GetAccord(accord.ID, userID)
	assert.Error(t, err)
}

func TestAccordHandler_AddTag(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord
	accord, err := service.CreateAccord(userID, &models.CreateAccordRequest{
		Name:            "Test",
		PyramidPosition: "top",
		VolumeMl:        10.0,
	})
	require.NoError(t, err)

	e := echo.New()
	reqBody := `{"tag": "newtag"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/accords/%s/tags", accord.ID), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id/tags")
	c.SetParamNames("id")
	c.SetParamValues(accord.ID)
	c.Set("userId", userID)

	err = handler.AddTag(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify tag was added
	fetched, err := service.GetAccord(accord.ID, userID)
	require.NoError(t, err)
	assert.Contains(t, fetched.Tags, "newtag")
}

func TestAccordHandler_RemoveTag(t *testing.T) {
	handler, service, tdb, userID := setupAccordHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Create an accord with tags
	accord, err := service.CreateAccord(userID, &models.CreateAccordRequest{
		Name:            "Test",
		PyramidPosition: "top",
		VolumeMl:        10.0,
		Tags:            []string{"tag1", "tag2"},
	})
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/accords/%s/tags/tag1", accord.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/accords/:id/tags/:tag")
	c.SetParamNames("id", "tag")
	c.SetParamValues(accord.ID, "tag1")
	c.Set("userId", userID)

	err = handler.RemoveTag(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify tag was removed
	fetched, err := service.GetAccord(accord.ID, userID)
	require.NoError(t, err)
	assert.NotContains(t, fetched.Tags, "tag1")
	assert.Contains(t, fetched.Tags, "tag2")
}
