package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type TagHandler struct {
	service *services.TagService
}

func NewTagHandler(service *services.TagService) *TagHandler {
	return &TagHandler{
		service: service,
	}
}

// GetAll retrieves all predefined tags
func (h *TagHandler) GetAll(c echo.Context) error {
	tags, err := h.service.GetAllTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to retrieve tags"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags": tags,
	})
}

// GetByCategory retrieves tags for a specific category
func (h *TagHandler) GetByCategory(c echo.Context) error {
	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Category is required"},
		})
	}

	tags, err := h.service.GetTagsByCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to retrieve tags"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags": tags,
	})
}

// Search searches for tags by partial match
func (h *TagHandler) Search(c echo.Context) error {
	search := c.QueryParam("q")

	tags, err := h.service.SearchTags(search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to search tags"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags": tags,
	})
}

// GetCategories retrieves all unique categories
func (h *TagHandler) GetCategories(c echo.Context) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to retrieve categories"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

// GetGrouped retrieves all tags grouped by category
func (h *TagHandler) GetGrouped(c echo.Context) error {
	grouped, err := h.service.GetTagsGroupedByCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to retrieve tags"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags": grouped,
	})
}
