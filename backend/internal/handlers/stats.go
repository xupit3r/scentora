package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type StatsHandler struct {
	accordService *services.AccordService
}

func NewStatsHandler(accordService *services.AccordService) *StatsHandler {
	return &StatsHandler{
		accordService: accordService,
	}
}

// GetStats returns statistics about the user's accord collection
func (h *StatsHandler) GetStats(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	stats, err := h.accordService.GetStatistics(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to retrieve statistics"},
		})
	}

	return c.JSON(http.StatusOK, stats)
}
