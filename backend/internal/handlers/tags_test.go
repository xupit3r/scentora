package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/services"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func setupTagHandler(t *testing.T) (*TagHandler, *testutil.TestDB) {
	tdb := testutil.SetupTestDB(t)
	tagRepo := repository.NewPredefinedTagRepository(tdb.DB)
	tagService := services.NewTagService(tagRepo)
	handler := NewTagHandler(tagService)
	return handler, tdb
}

func TestTagHandler_GetAll(t *testing.T) {
	handler, tdb := setupTagHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAll(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "tags")
}

func TestTagHandler_GetByCategory(t *testing.T) {
	handler, tdb := setupTagHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tags/category/scent_family", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tags/category/:category")
	c.SetParamNames("category")
	c.SetParamValues("scent_family")

	err := handler.GetByCategory(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "tags")
}

func TestTagHandler_Search(t *testing.T) {
	handler, tdb := setupTagHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tags/search?q=floral", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Search(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "tags")
}

func TestTagHandler_GetCategories(t *testing.T) {
	handler, tdb := setupTagHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tags/categories", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCategories(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "categories")
}

func TestTagHandler_GetGrouped(t *testing.T) {
	handler, tdb := setupTagHandler(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/tags/grouped", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetGrouped(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "tags")
}
