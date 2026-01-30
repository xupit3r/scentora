package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type PerfumeHandler struct {
	service   *services.PerfumeService
	validator *validator.Validate
}

func NewPerfumeHandler(service *services.PerfumeService) *PerfumeHandler {
	return &PerfumeHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *PerfumeHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	var req models.CreatePerfumeRequest
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

	result, err := h.service.Create(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *PerfumeHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result, err := h.service.Get(id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PerfumeHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)

	filters := map[string]string{
		"search":        c.QueryParam("search"),
		"year":          c.QueryParam("year"),
		"concentration": c.QueryParam("concentration"),
		"note":          c.QueryParam("note"),
	}

	results, err := h.service.List(userID, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, results)
}

func (h *PerfumeHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req models.UpdatePerfumeRequest
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

	result, err := h.service.Update(id, userID, &req)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PerfumeHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.service.Delete(id, userID); err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.NoContent(http.StatusNoContent)
}
