package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestRecipeNoteRepository_Create(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, version := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	note := &models.RecipeNote{
		RecipeID:  recipe.ID,
		VersionID: &version.ID,
		Content:   "This version smells great!",
		NoteType:  "testing",
	}

	err := repo.Create(note)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	if note.ID == "" {
		t.Error("Expected note ID to be set")
	}
}

func TestRecipeNoteRepository_FindByRecipeID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, version := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	
	// Create general note
	note1 := &models.RecipeNote{
		RecipeID: recipe.ID,
		Content:  "General recipe note",
		NoteType: "general",
	}
	err := repo.Create(note1)
	if err != nil {
		t.Fatalf("Failed to create note 1: %v", err)
	}

	// Create version-specific note
	note2 := &models.RecipeNote{
		RecipeID:  recipe.ID,
		VersionID: &version.ID,
		Content:   "Version-specific note",
		NoteType:  "testing",
	}
	err = repo.Create(note2)
	if err != nil {
		t.Fatalf("Failed to create note 2: %v", err)
	}

	notes, err := repo.FindByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to find notes: %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(notes))
	}
}

func TestRecipeNoteRepository_FindByVersionID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, version := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	
	note := &models.RecipeNote{
		RecipeID:  recipe.ID,
		VersionID: &version.ID,
		Content:   "Version note",
		NoteType:  "testing",
	}
	err := repo.Create(note)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	notes, err := repo.FindByVersionID(version.ID)
	if err != nil {
		t.Fatalf("Failed to find notes: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}
	if notes[0].Content != "Version note" {
		t.Error("Expected correct note content")
	}
}

func TestRecipeNoteRepository_FindByRecipeIDAndType(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	
	// Create notes with different types
	types := []string{"general", "testing", "observation"}
	for _, noteType := range types {
		note := &models.RecipeNote{
			RecipeID: recipe.ID,
			Content:  "Note of type " + noteType,
			NoteType: noteType,
		}
		err := repo.Create(note)
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
	}

	testingNotes, err := repo.FindByRecipeIDAndType(recipe.ID, "testing")
	if err != nil {
		t.Fatalf("Failed to find testing notes: %v", err)
	}

	if len(testingNotes) != 1 {
		t.Errorf("Expected 1 testing note, got %d", len(testingNotes))
	}
	if testingNotes[0].NoteType != "testing" {
		t.Error("Expected note type 'testing'")
	}
}

func TestRecipeNoteRepository_Update(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	note := &models.RecipeNote{
		RecipeID: recipe.ID,
		Content:  "Original content",
		NoteType: "general",
	}
	err := repo.Create(note)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	note.Content = "Updated content"
	note.NoteType = "observation"

	err = repo.Update(note)
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}

	found, err := repo.FindByID(note.ID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}

	if found.Content != "Updated content" {
		t.Error("Expected updated content")
	}
	if found.NoteType != "observation" {
		t.Error("Expected updated note type")
	}
}

func TestRecipeNoteRepository_Delete(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	note := &models.RecipeNote{
		RecipeID: recipe.ID,
		Content:  "To delete",
		NoteType: "general",
	}
	err := repo.Create(note)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	err = repo.Delete(note.ID)
	if err != nil {
		t.Fatalf("Failed to delete note: %v", err)
	}

	_, err = repo.FindByID(note.ID)
	if err == nil {
		t.Error("Expected error finding deleted note")
	}
}

func TestRecipeNoteRepository_CountByRecipeID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	user, recipe, _ := createTestRecipeWithVersion(t, tdb)
	_ = user

	repo := NewRecipeNoteRepository(tdb.DB)
	
	// Create 3 notes
	for i := 0; i < 3; i++ {
		note := &models.RecipeNote{
			RecipeID: recipe.ID,
			Content:  testutil.UniqueString("Note"),
			NoteType: "general",
		}
		err := repo.Create(note)
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
	}

	count, err := repo.CountByRecipeID(recipe.ID)
	if err != nil {
		t.Fatalf("Failed to count notes: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}
