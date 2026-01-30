package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type InvitationHandler struct {
	service   *services.InvitationService
	validator *validator.Validate
}

func NewInvitationHandler(service *services.InvitationService) *InvitationHandler {
	return &InvitationHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *InvitationHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	var req models.CreateInvitationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid request body"},
		})
	}

	// Set default expiration if not provided
	if req.ExpiresInDays == 0 {
		req.ExpiresInDays = 7
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Validation failed", Details: err.Error()},
		})
	}

	invitation, err := h.service.Create(userID, req.Email, req.ExpiresInDays)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"invitation": invitation,
	})
}

func (h *InvitationHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	invitations, err := h.service.List(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"invitations": invitations,
	})
}

func (h *InvitationHandler) Revoke(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invitation code is required"},
		})
	}

	err := h.service.Revoke(code, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Invitation revoked successfully",
	})
}
