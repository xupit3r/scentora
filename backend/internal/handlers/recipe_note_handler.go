package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/services"
)

type RecipeNoteHandler struct {
	noteService *services.RecipeNoteService
}

func NewRecipeNoteHandler(noteService *services.RecipeNoteService) *RecipeNoteHandler {
	return &RecipeNoteHandler{
		noteService: noteService,
	}
}

// CreateNote handles POST /api/recipes/:id/notes
func (h *RecipeNoteHandler) CreateNote(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	var req models.CreateRecipeNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	note, err := h.noteService.CreateNote(recipeID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, note)
}

// ListNotes handles GET /api/recipes/:id/notes
func (h *RecipeNoteHandler) ListNotes(c echo.Context) error {
	userID := c.Get("userId").(string)
	recipeID := c.Param("id")

	notes, err := h.noteService.GetNotes(recipeID, userID, nil)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, notes)
}

// UpdateNote handles PUT /api/recipes/:id/notes/:noteId
func (h *RecipeNoteHandler) UpdateNote(c echo.Context) error {
	userID := c.Get("userId").(string)
	noteID := c.Param("noteId")

	var req models.UpdateRecipeNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	note, err := h.noteService.UpdateNote(noteID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, note)
}

// DeleteNote handles DELETE /api/recipes/:id/notes/:noteId
func (h *RecipeNoteHandler) DeleteNote(c echo.Context) error {
	userID := c.Get("userId").(string)
	noteID := c.Param("noteId")

	if err := h.noteService.DeleteNote(noteID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Note deleted successfully"})
}
