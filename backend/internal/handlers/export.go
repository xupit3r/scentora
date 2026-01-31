package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type ExportHandler struct {
	accordService *services.AccordService
}

func NewExportHandler(accordService *services.AccordService) *ExportHandler {
	return &ExportHandler{
		accordService: accordService,
	}
}

// Export exports all accords for a user as JSON
func (h *ExportHandler) Export(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	accords, err := h.accordService.ListAccords(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Failed to export accords"},
		})
	}

	// Convert to AccordResponse format
	accordResponses := make([]*models.AccordResponse, len(accords))
	for i, accord := range accords {
		accordResponses[i] = &models.AccordResponse{
			ID:                 accord.ID,
			UserID:             accord.UserID,
			Name:               accord.Name,
			PyramidPosition:    accord.PyramidPosition,
			VolumeMl:           accord.VolumeMl,
			VolumeDrops:        accord.VolumeDrops,
			Supplier:           accord.Supplier,
			PurchaseDate:       accord.PurchaseDate,
			DilutionPercentage: accord.DilutionPercentage,
			Notes:              accord.Notes,
			Tags:               accord.Tags,
			CreatedAt:          accord.CreatedAt,
			UpdatedAt:          accord.UpdatedAt,
		}
	}

	exportData := models.ExportData{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Accords:    accordResponses,
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=scentora-export.json")
	return c.JSON(http.StatusOK, exportData)
}

// Import imports accords from JSON
func (h *ExportHandler) Import(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Unauthorized"},
		})
	}

	var importData models.ImportData
	if err := c.Bind(&importData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorDetail{Message: "Invalid import data"},
		})
	}

	// Import accords
	imported := 0
	failed := 0
	var errors []string

	for _, accord := range importData.Accords {
		// Create new accord from import
		createReq := models.CreateAccordRequest{
			Name:               accord.Name,
			PyramidPosition:    accord.PyramidPosition,
			VolumeMl:           accord.VolumeMl,
			Supplier:           accord.Supplier,
			PurchaseDate:       accord.PurchaseDate,
			DilutionPercentage: accord.DilutionPercentage,
			Notes:              accord.Notes,
			Tags:               accord.Tags,
		}

		_, err := h.accordService.CreateAccord(userID, &createReq)
		if err != nil {
			failed++
			errors = append(errors, err.Error())
		} else {
			imported++
		}
	}

	result := models.ImportResult2{
		TotalRecords:    len(importData.Accords),
		ImportedRecords: imported,
		FailedRecords:   failed,
		Errors:          errors,
	}

	return c.JSON(http.StatusOK, result)
}

