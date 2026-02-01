package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeTagHandler struct {
	tagService *services.RecipeTagService
}

func NewRecipeTagHandler(tagService *services.RecipeTagService) *RecipeTagHandler {
	return &RecipeTagHandler{
		tagService: tagService,
	}
}

// AddTag handles POST /api/recipes/:id/tags
func (h *RecipeTagHandler) AddTag(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	var req struct {
		Tag string `json:"tag" binding:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.tagService.AddTag(recipeID, userID, req.Tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Tag added successfully", "tag": req.Tag})
}

// RemoveTag handles DELETE /api/recipes/:id/tags/:tag
func (h *RecipeTagHandler) RemoveTag(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")
	tag := c.Param("tag")

	if err := h.tagService.RemoveTag(recipeID, userID, tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Tag removed successfully"})
}

// GetPopularTags handles GET /api/recipes/tags/popular
func (h *RecipeTagHandler) GetPopularTags(c echo.Context) error {
	userID := c.Get("userId").(string)

	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	tags, err := h.tagService.GetPopularTags(userID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tags)
}
