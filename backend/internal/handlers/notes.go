package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type NotesHandler struct {
	perfumeRepo *repository.PerfumeRepository
}

func NewNotesHandler(perfumeRepo *repository.PerfumeRepository) *NotesHandler {
	return &NotesHandler{perfumeRepo: perfumeRepo}
}

func (h *NotesHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	notes, err := h.perfumeRepo.GetAllNotes(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, notes)
}
