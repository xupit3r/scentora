package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type AccordHandler struct {
	service   *services.AccordService
	validator *validator.Validate
}

func NewAccordHandler(service *services.AccordService) *AccordHandler {
	return &AccordHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create creates a new accord
func (h *AccordHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	var req models.CreateAccordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid request body"},
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Validation failed", Details: err.Error()},
		})
	}

	accord, err := h.service.CreateAccord(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"accord": accord,
	})
}

// Get retrieves a single accord
func (h *AccordHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accordID := c.Param("id")
	if accordID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord ID is required"},
		})
	}

	accord, err := h.service.GetAccord(accordID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord not found"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accord": accord,
	})
}

// List retrieves all accords for the authenticated user
func (h *AccordHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	// Check for search/filter parameters
	position := c.QueryParam("position")
	supplier := c.QueryParam("supplier")
	search := c.QueryParam("search")
	tags := c.QueryParams()["tags"]

	var minVolume, maxVolume *float64
	if minStr := c.QueryParam("minVolume"); minStr != "" {
		if val, err := strconv.ParseFloat(minStr, 64); err == nil {
			minVolume = &val
		}
	}
	if maxStr := c.QueryParam("maxVolume"); maxStr != "" {
		if val, err := strconv.ParseFloat(maxStr, 64); err == nil {
			maxVolume = &val
		}
	}

	var posPtr, supplierPtr, searchPtr *string
	if position != "" {
		posPtr = &position
	}
	if supplier != "" {
		supplierPtr = &supplier
	}
	if search != "" {
		searchPtr = &search
	}

	// Use search if any filters are provided
	if posPtr != nil || minVolume != nil || maxVolume != nil || supplierPtr != nil || searchPtr != nil || len(tags) > 0 {
		accords, err := h.service.SearchAccords(userID, posPtr, minVolume, maxVolume, supplierPtr, searchPtr, tags)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: models.ErrorDetail{Message: "Failed to search accords"},
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"accords": accords,
		})
	}

	// Otherwise list all
	accords, err := h.service.ListAccords(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to list accords"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accords": accords,
	})
}

// Update updates an existing accord
func (h *AccordHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accordID := c.Param("id")
	if accordID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord ID is required"},
		})
	}

	var req models.UpdateAccordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid request body"},
		})
	}

	accord, err := h.service.UpdateAccord(accordID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accord": accord,
	})
}

// Delete deletes an accord
func (h *AccordHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accordID := c.Param("id")
	if accordID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord ID is required"},
		})
	}

	err := h.service.DeleteAccord(accordID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Accord deleted successfully",
	})
}

// AddTag adds a tag to an accord
func (h *AccordHandler) AddTag(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accordID := c.Param("id")
	if accordID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord ID is required"},
		})
	}

	var req models.AddTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid request body"},
		})
	}

	if req.Tag == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Tag is required"},
		})
	}

	err := h.service.AddTagToAccord(accordID, userID, req.Tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tag added successfully",
	})
}

// RemoveTag removes a tag from an accord
func (h *AccordHandler) RemoveTag(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accordID := c.Param("id")
	tag := c.Param("tag")

	if accordID == "" || tag == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Accord ID and tag are required"},
		})
	}

	err := h.service.RemoveTagFromAccord(accordID, userID, tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tag removed successfully",
	})
}
