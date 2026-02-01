package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeCollectionHandler struct {
	collectionService *services.RecipeCollectionService
}

func NewRecipeCollectionHandler(collectionService *services.RecipeCollectionService) *RecipeCollectionHandler {
	return &RecipeCollectionHandler{
		collectionService: collectionService,
	}
}

// CreateCollection handles POST /api/collections
func (h *RecipeCollectionHandler) CreateCollection(c echo.Context) error {
	userID := c.Get("userId").(string)

	var req models.CreateRecipeCollectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	collection, err := h.collectionService.CreateCollection(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, collection)
}

// ListCollections handles GET /api/collections
func (h *RecipeCollectionHandler) ListCollections(c echo.Context) error {
	userID := c.Get("userId").(string)

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

	collections, err := h.collectionService.ListCollections(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, collections)
}

// GetCollection handles GET /api/collections/:id
func (h *RecipeCollectionHandler) GetCollection(c echo.Context) error {
	userID := c.Get("userId").(string)
	collectionID := c.Param("id")

	collection, err := h.collectionService.GetCollection(collectionID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, collection)
}

// UpdateCollection handles PUT /api/collections/:id
func (h *RecipeCollectionHandler) UpdateCollection(c echo.Context) error {
	userID := c.Get("userId").(string)
	collectionID := c.Param("id")

	var req models.UpdateRecipeCollectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	collection, err := h.collectionService.UpdateCollection(collectionID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, collection)
}

// DeleteCollection handles DELETE /api/collections/:id
func (h *RecipeCollectionHandler) DeleteCollection(c echo.Context) error {
	userID := c.Get("userId").(string)
	collectionID := c.Param("id")

	if err := h.collectionService.DeleteCollection(collectionID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Collection deleted successfully"})
}

// AddRecipeToCollection handles POST /api/collections/:id/recipes
func (h *RecipeCollectionHandler) AddRecipeToCollection(c echo.Context) error {
	userID := c.Get("userId").(string)
	collectionID := c.Param("id")

	var req struct {
		RecipeID string `json:"recipeId" binding:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.collectionService.AddRecipeToCollection(collectionID, req.RecipeID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Recipe added to collection"})
}

// RemoveRecipeFromCollection handles DELETE /api/collections/:id/recipes/:recipeId
func (h *RecipeCollectionHandler) RemoveRecipeFromCollection(c echo.Context) error {
	userID := c.Get("userId").(string)
	collectionID := c.Param("id")
	recipeID := c.Param("recipeId")

	if err := h.collectionService.RemoveRecipeFromCollection(collectionID, recipeID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Recipe removed from collection"})
}
