package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeIngredientHandler struct {
	ingredientService *services.RecipeIngredientService
}

func NewRecipeIngredientHandler(ingredientService *services.RecipeIngredientService) *RecipeIngredientHandler {
	return &RecipeIngredientHandler{
		ingredientService: ingredientService,
	}
}

// AddIngredient handles POST /api/recipes/:id/versions/:versionId/ingredients
func (h *RecipeIngredientHandler) AddIngredient(c echo.Context) error {
	userID := c.Get("userId").(string)
	versionID := c.Param("versionId")

	var req models.CreateRecipeIngredientRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ingredient, err := h.ingredientService.AddIngredient(versionID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, ingredient)
}

// UpdateIngredient handles PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId
func (h *RecipeIngredientHandler) UpdateIngredient(c echo.Context) error {
	userID := c.Get("userId").(string)
	ingredientID := c.Param("ingredientId")

	var req models.UpdateRecipeIngredientRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ingredient, err := h.ingredientService.UpdateIngredient(ingredientID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ingredient)
}

// RemoveIngredient handles DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId
func (h *RecipeIngredientHandler) RemoveIngredient(c echo.Context) error {
	userID := c.Get("userId").(string)
	ingredientID := c.Param("ingredientId")

	if err := h.ingredientService.DeleteIngredient(ingredientID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Ingredient removed successfully"})
}
