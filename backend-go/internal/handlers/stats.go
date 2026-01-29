package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type StatsHandler struct {
	perfumeRepo *repository.PerfumeRepository
	journalRepo *repository.JournalRepository
}

func NewStatsHandler(perfumeRepo *repository.PerfumeRepository, journalRepo *repository.JournalRepository) *StatsHandler {
	return &StatsHandler{
		perfumeRepo: perfumeRepo,
		journalRepo: journalRepo,
	}
}

func (h *StatsHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)

	perfumeCount, err := h.perfumeRepo.Count(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	journalCount, err := h.journalRepo.Count(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	stats := map[string]interface{}{
		"perfumes": perfumeCount,
		"entries":  journalCount,
	}

	return c.JSON(http.StatusOK, stats)
}
