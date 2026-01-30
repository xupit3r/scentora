package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type ExportHandler struct {
	perfumeRepo *repository.PerfumeRepository
	journalRepo *repository.JournalRepository
	validator   *validator.Validate
}

func NewExportHandler(perfumeRepo *repository.PerfumeRepository, journalRepo *repository.JournalRepository) *ExportHandler {
	return &ExportHandler{
		perfumeRepo: perfumeRepo,
		journalRepo: journalRepo,
		validator:   validator.New(),
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

	export := &models.ExportResponse{
		Version:        "1.0",
		ExportDate:     time.Now().Format(time.RFC3339),
		Perfumes:       perfumeResponses,
		JournalEntries: journals,
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=\"scentora-collection.json\"")
	c.Response().Header().Set("Content-Type", "application/json")

	return c.JSON(http.StatusOK, export)
}

func (h *ExportHandler) Import(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req models.ImportCollectionRequest
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

	if req.Perfumes == nil || len(req.Perfumes) == 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid import data format"},
		})
	}

	results := &models.ImportResult{
		PerfumesImported:       0,
		JournalEntriesImported: 0,
		Errors:                 []string{},
	}

	// Import perfumes
	for _, importPerfume := range req.Perfumes {
		perfume := &models.Perfume{
			UserID:        userID,
			Name:          importPerfume.Name,
			Designer:      importPerfume.Designer,
			Year:          importPerfume.Year,
			Concentration: importPerfume.Concentration,
			TopNotes:      pq.StringArray(importPerfume.Pyramid.Top),
			MiddleNotes:   pq.StringArray(importPerfume.Pyramid.Middle),
			BaseNotes:     pq.StringArray(importPerfume.Pyramid.Base),
			Description:   importPerfume.Description,
			ImageURL:      importPerfume.ImageURL,
		}

		if err := h.perfumeRepo.Create(perfume); err != nil {
			results.Errors = append(results.Errors, fmt.Sprintf("Failed to import perfume: %s", importPerfume.Name))
		} else {
			results.PerfumesImported++
		}
	}

	// Import journal entries (if present)
	if req.JournalEntries != nil {
		for _, importJournal := range req.JournalEntries {
			// Note: perfumeID from import may not exist, so we skip invalid entries
			entry := &models.JournalEntry{
				UserID:    userID,
				PerfumeID: importJournal.PerfumeID,
				Date:      importJournal.Date,
				Content:   importJournal.Content,
				Rating:    importJournal.Rating,
				Occasion:  importJournal.Occasion,
				Weather:   importJournal.Weather,
			}

			if err := h.journalRepo.Create(entry); err != nil {
				results.Errors = append(results.Errors, "Failed to import journal entry")
			} else {
				results.JournalEntriesImported++
			}
		}
	}

	return c.JSON(http.StatusOK, results)
}
