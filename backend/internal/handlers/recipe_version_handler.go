package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeVersionHandler struct {
	versionService *services.RecipeVersionService
}

func NewRecipeVersionHandler(versionService *services.RecipeVersionService) *RecipeVersionHandler {
	return &RecipeVersionHandler{
		versionService: versionService,
	}
}

// CreateVersion handles POST /api/recipes/:id/versions
func (h *RecipeVersionHandler) CreateVersion(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	var req models.CreateRecipeVersionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	version, err := h.versionService.CreateVersion(recipeID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, version)
}

// ListVersions handles GET /api/recipes/:id/versions
func (h *RecipeVersionHandler) ListVersions(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	versions, err := h.versionService.ListVersions(recipeID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, versions)
}

// GetVersion handles GET /api/recipes/:id/versions/:versionId
func (h *RecipeVersionHandler) GetVersion(c echo.Context) error {
	versionID := c.Param("versionId")

	version, err := h.versionService.GetVersion(versionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, version)
}

// ActivateVersion handles POST /api/recipes/:id/versions/:versionNumber/activate
func (h *RecipeVersionHandler) ActivateVersion(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")
	versionNumberStr := c.Param("versionNumber")

	versionNumber, err := strconv.Atoi(versionNumberStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid version number"})
	}

	// Get all versions to find the one with this version number
	versions, err := h.versionService.ListVersions(recipeID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var versionID string
	for _, v := range versions {
		if v.VersionNumber == versionNumber {
			versionID = v.ID
			break
		}
	}

	if versionID == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Version not found"})
	}

	if err := h.versionService.SetActiveVersion(recipeID, versionID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Version activated successfully"})
}

// DuplicateVersion handles POST /api/recipes/:id/versions/:versionId/duplicate
func (h *RecipeVersionHandler) DuplicateVersion(c echo.Context) error {
	userID := c.Get("userId").(string)
	versionID := c.Param("versionId")

	newVersion, err := h.versionService.DuplicateVersion(versionID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newVersion)
}
