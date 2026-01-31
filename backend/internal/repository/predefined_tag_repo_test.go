package repository

import (
	"testing"

	"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestPredefinedTagRepository_GetAll(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)
	tags, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all tags: %v", err)
	}

	// Should have 57 predefined tags (from seeding)
	if len(tags) != 57 {
		t.Errorf("Expected 57 predefined tags, got %d", len(tags))
	}
}

func TestPredefinedTagRepository_GetByCategory(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	// Test known category
	tags, err := repo.GetByCategory("scent_family")
	if err != nil {
		t.Fatalf("Failed to get tags by category: %v", err)
	}

	if len(tags) == 0 {
		t.Error("Expected at least one tag in scent_family category")
	}

	// Verify all tags have correct category
	for _, tag := range tags {
		if tag.Category != "scent_family" {
			t.Errorf("Expected category scent_family, got %s", tag.Category)
		}
	}
}

func TestPredefinedTagRepository_GetByCategory_Empty(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	tags, err := repo.GetByCategory("nonexistent_category")
	if err != nil {
		t.Fatalf("Failed to get tags by category: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags for nonexistent category, got %d", len(tags))
	}
}

func TestPredefinedTagRepository_Search(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	// Search for tags starting with "fresh"
	tags, err := repo.Search("fresh")
	if err != nil {
		t.Fatalf("Failed to search tags: %v", err)
	}

	// Should find at least one tag
	if len(tags) == 0 {
		t.Error("Expected at least one tag matching 'fresh'")
	}

	// Verify all tags start with "fresh"
	for _, tag := range tags {
		if len(tag.Tag) < 5 || tag.Tag[:5] != "fresh" {
			t.Errorf("Expected tag to start with 'fresh', got %s", tag.Tag)
		}
	}
}

func TestPredefinedTagRepository_Search_CaseInsensitive(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	// Search with uppercase should work (ILIKE)
	tags, err := repo.Search("FRESH")
	if err != nil {
		t.Fatalf("Failed to search tags: %v", err)
	}

	if len(tags) == 0 {
		t.Error("Expected case-insensitive search to find results")
	}
}

func TestPredefinedTagRepository_Search_NoResults(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	tags, err := repo.Search("xyz123nonexistent")
	if err != nil {
		t.Fatalf("Failed to search tags: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags for nonexistent search, got %d", len(tags))
	}
}

func TestPredefinedTagRepository_GetAllCategories(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Teardown(t)

	repo := NewPredefinedTagRepository(tdb.DB)

	categories, err := repo.GetAllCategories()
	if err != nil {
		t.Fatalf("Failed to get all categories: %v", err)
	}

	// Should have 9 categories (from seeding)
	if len(categories) != 9 {
		t.Errorf("Expected 9 categories, got %d", len(categories))
	}

	expectedCategories := map[string]bool{
		"scent_family": true,
		"character":    true,
		"mood":         true,
		"season":       true,
		"occasion":     true,
		"time_of_day":  true,
		"longevity":    true,
		"sillage":      true,
		"ingredients":  true,
	}

	for _, cat := range categories {
		if !expectedCategories[cat] {
			t.Errorf("Unexpected category: %s", cat)
		}
	}
}
