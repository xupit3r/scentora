package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeHandler struct {
	recipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

// CreateRecipe handles POST /api/recipes
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	userID := c.Get("userId").(string)

	var req models.CreateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	recipe, err := h.recipeService.CreateRecipe(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, recipe)
}

// GetRecipe handles GET /api/recipes/:id
func (h *RecipeHandler) GetRecipe(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	recipe, err := h.recipeService.GetRecipe(recipeID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, recipe)
}

// ListRecipes handles GET /api/recipes
func (h *RecipeHandler) ListRecipes(c echo.Context) error {
	userID := c.Get("userId").(string)

	// Get query parameters
	status := c.QueryParam("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	limit := 100
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	recipes, err := h.recipeService.ListRecipes(userID, statusPtr, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, recipes)
}

// UpdateRecipe handles PUT /api/recipes/:id
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	var req models.UpdateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	recipe, err := h.recipeService.UpdateRecipe(recipeID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe handles DELETE /api/recipes/:id
func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	if err := h.recipeService.DeleteRecipe(recipeID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Recipe deleted successfully"})
}

// SearchRecipes handles GET /api/recipes/search
func (h *RecipeHandler) SearchRecipes(c echo.Context) error {
	userID := c.Get("userId").(string)

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query required"})
	}

	limit := 100
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	recipes, err := h.recipeService.SearchRecipes(userID, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, recipes)
}
