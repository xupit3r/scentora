package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type RecipeNoteRepository struct {
	db *sqlx.DB
}

func NewRecipeNoteRepository(db *sqlx.DB) *RecipeNoteRepository {
	return &RecipeNoteRepository{db: db}
}

// Create creates a new recipe note
func (r *RecipeNoteRepository) Create(note *models.RecipeNote) error {
	query := `
		INSERT INTO recipe_notes (
			recipe_id, version_id, content, note_type, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		note.RecipeID,
		note.VersionID,
		note.Content,
		note.NoteType,
	).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

// FindByID retrieves a note by ID
func (r *RecipeNoteRepository) FindByID(id string) (*models.RecipeNote, error) {
	var note models.RecipeNote
	query := `
		SELECT id, recipe_id, version_id, content, note_type, created_at, updated_at
		FROM recipe_notes
		WHERE id = $1
	`
	err := r.db.Get(&note, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found")
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// FindByRecipeID retrieves all notes for a recipe
func (r *RecipeNoteRepository) FindByRecipeID(recipeID string) ([]*models.RecipeNote, error) {
	var notes []*models.RecipeNote
	query := `
		SELECT id, recipe_id, version_id, content, note_type, created_at, updated_at
		FROM recipe_notes
		WHERE recipe_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&notes, query, recipeID)
	return notes, err
}

// FindByVersionID retrieves all notes for a specific version
func (r *RecipeNoteRepository) FindByVersionID(versionID string) ([]*models.RecipeNote, error) {
	var notes []*models.RecipeNote
	query := `
		SELECT id, recipe_id, version_id, content, note_type, created_at, updated_at
		FROM recipe_notes
		WHERE version_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&notes, query, versionID)
	return notes, err
}

// FindByRecipeIDAndType retrieves notes filtered by type
func (r *RecipeNoteRepository) FindByRecipeIDAndType(recipeID, noteType string) ([]*models.RecipeNote, error) {
	var notes []*models.RecipeNote
	query := `
		SELECT id, recipe_id, version_id, content, note_type, created_at, updated_at
		FROM recipe_notes
		WHERE recipe_id = $1 AND note_type = $2
		ORDER BY created_at DESC
	`
	err := r.db.Select(&notes, query, recipeID, noteType)
	return notes, err
}

// Update updates a note's content and type
func (r *RecipeNoteRepository) Update(note *models.RecipeNote) error {
	query := `
		UPDATE recipe_notes
		SET content = $1, note_type = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`
	err := r.db.QueryRow(query, note.Content, note.NoteType, note.ID).Scan(&note.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("note not found")
	}
	return err
}

// Delete deletes a note
func (r *RecipeNoteRepository) Delete(id string) error {
	query := `DELETE FROM recipe_notes WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("note not found")
	}
	return nil
}

// CountByRecipeID returns the number of notes for a recipe
func (r *RecipeNoteRepository) CountByRecipeID(recipeID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM recipe_notes WHERE recipe_id = $1`
	err := r.db.Get(&count, query, recipeID)
	return count, err
}
