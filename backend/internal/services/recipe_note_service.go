package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type RecipeNoteService struct {
	noteRepo   *repository.RecipeNoteRepository
	recipeRepo *repository.RecipeRepository
}

func NewRecipeNoteService(
	noteRepo *repository.RecipeNoteRepository,
	recipeRepo *repository.RecipeRepository,
) *RecipeNoteService {
	return &RecipeNoteService{
		noteRepo:   noteRepo,
		recipeRepo: recipeRepo,
	}
}

// CreateNote creates a new note for a recipe
func (s *RecipeNoteService) CreateNote(recipeID, userID string, req *models.CreateRecipeNoteRequest) (*models.RecipeNote, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Validate input
	if req.Content == "" {
		return nil, errors.New("note content is required")
	}

	// Validate note type
	if req.NoteType != nil {
		validTypes := map[string]bool{"general": true, "version": true, "test": true}
		if !validTypes[*req.NoteType] {
			return nil, errors.New("note type must be one of: general, version, test")
		}
	}

	// Create note
	noteType := "general" // default
	if req.NoteType != nil {
		noteType = *req.NoteType
	}
	
	note := &models.RecipeNote{
		RecipeID: recipeID,
		Content:  req.Content,
		NoteType: noteType,
	}

	err = s.noteRepo.Create(note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

// GetNotes retrieves notes for a recipe
func (s *RecipeNoteService) GetNotes(recipeID, userID string, noteType *string) ([]*models.RecipeNote, error) {
	// Verify recipe exists and belongs to user
	_, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	var notes []*models.RecipeNote
	if noteType != nil && *noteType != "" {
		notes, err = s.noteRepo.FindByRecipeIDAndType(recipeID, *noteType)
	} else {
		notes, err = s.noteRepo.FindByRecipeID(recipeID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}

	return notes, nil
}

// UpdateNote updates a note
func (s *RecipeNoteService) UpdateNote(noteID, userID string, req *models.UpdateRecipeNoteRequest) (*models.RecipeNote, error) {
	// Get note
	note, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(userID, note.RecipeID)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}

	// Validate content
	if req.Content == "" {
		return nil, errors.New("note content is required")
	}

	// Update note
	note.Content = req.Content
	err = s.noteRepo.Update(note)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return note, nil
}

// DeleteNote deletes a note
func (s *RecipeNoteService) DeleteNote(noteID, userID string) error {
	// Get note
	note, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Verify recipe exists and belongs to user
	_, err = s.recipeRepo.FindByID(userID, note.RecipeID)
	if err != nil {
		return fmt.Errorf("recipe not found: %w", err)
	}

	// Delete note
	err = s.noteRepo.Delete(noteID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
