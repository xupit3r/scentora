package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type AuthHandler struct {
	service   *services.AuthService
	validator *validator.Validate
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
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

	result, err := h.service.Register(req.Email, req.Username, req.Password, req.InvitationCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
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

	result, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req models.RefreshRequest
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

	result, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var req models.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid request body"},
		})
	}

	if err := h.service.Logout(req.RefreshToken); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	if err := h.service.LogoutAll(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out from all devices",
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "User not found"},
		})
	}

	return c.JSON(http.StatusOK, user)
}
