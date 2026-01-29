package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type ExportHandler struct {
	perfumeRepo *repository.PerfumeRepository
	journalRepo *repository.JournalRepository
}

func NewExportHandler(perfumeRepo *repository.PerfumeRepository, journalRepo *repository.JournalRepository) *ExportHandler {
	return &ExportHandler{
		perfumeRepo: perfumeRepo,
		journalRepo: journalRepo,
	}
}

func (h *ExportHandler) Export(c echo.Context) error {
	userID := middleware.GetUserID(c)

	perfumes, err := h.perfumeRepo.List(userID, map[string]string{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	journals, err := h.journalRepo.ListAll(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	// Convert perfumes to responses
	perfumeResponses := make([]*models.PerfumeResponse, len(perfumes))
	for i, p := range perfumes {
		perfumeResponses[i] = p.ToResponse()
	}

	export := map[string]interface{}{
		"perfumes": perfumeResponses,
		"journal":  journals,
	}

	return c.JSON(http.StatusOK, export)
}
