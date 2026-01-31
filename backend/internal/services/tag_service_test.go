package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/testutil"
)

func setupTagService(t *testing.T) (*TagService, *testutil.TestDB) {
	tdb := testutil.SetupTestDB(t)
	tagRepo := repository.NewPredefinedTagRepository(tdb.DB)
	service := NewTagService(tagRepo)
	return service, tdb
}

func TestTagService_GetAllTags(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get all tags (should return seeded tags from migrations)
	tags, err := service.GetAllTags()

	require.NoError(t, err)
	assert.NotEmpty(t, tags, "Should have seeded tags")
	assert.Greater(t, len(tags), 50, "Should have at least 50 predefined tags")
}

func TestTagService_GetTagsByCategory(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	testCases := []struct {
		category     string
		expectedTags []string
	}{
		{
			category:     "scent_family",
			expectedTags: []string{"citrus", "floral", "woody", "oriental"},
		},
		{
			category:     "character",
			expectedTags: []string{"fresh", "warm", "sweet", "spicy"},
		},
		{
			category:     "season",
			expectedTags: []string{"spring", "summer", "fall", "winter"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.category, func(t *testing.T) {
			tags, err := service.GetTagsByCategory(tc.category)
			require.NoError(t, err)
			assert.NotEmpty(t, tags)

			// Verify expected tags are present
			tagStrings := make([]string, len(tags))
			for i, tag := range tags {
				tagStrings[i] = tag.Tag
				assert.Equal(t, tc.category, tag.Category)
			}

			for _, expectedTag := range tc.expectedTags {
				assert.Contains(t, tagStrings, expectedTag,
					"Category %s should contain tag %s", tc.category, expectedTag)
			}
		})
	}
}

func TestTagService_GetTagsByCategory_Empty(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get tags for non-existent category
	tags, err := service.GetTagsByCategory("nonexistent_category")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagService_SearchTags(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	testCases := []struct {
		search       string
		shouldContain []string
	}{
		{
			search:       "flo",
			shouldContain: []string{"floral"},
		},
		{
			search:       "wood",
			shouldContain: []string{"woody"},
		},
		{
			search:       "fresh",
			shouldContain: []string{"fresh"},
		},
		{
			search:       "sum",
			shouldContain: []string{"summer"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.search, func(t *testing.T) {
			tags, err := service.SearchTags(tc.search)
			require.NoError(t, err)

			tagStrings := make([]string, len(tags))
			for i, tag := range tags {
				tagStrings[i] = tag.Tag
			}

			for _, expected := range tc.shouldContain {
				assert.Contains(t, tagStrings, expected,
					"Search '%s' should find tag '%s'", tc.search, expected)
			}
		})
	}
}

func TestTagService_SearchTags_CaseInsensitive(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Search with different cases
	lowerResults, err := service.SearchTags("floral")
	require.NoError(t, err)

	upperResults, err := service.SearchTags("FLORAL")
	require.NoError(t, err)

	mixedResults, err := service.SearchTags("FlOrAl")
	require.NoError(t, err)

	// All should return the same results
	assert.Equal(t, len(lowerResults), len(upperResults))
	assert.Equal(t, len(lowerResults), len(mixedResults))
}

func TestTagService_SearchTags_EmptyString(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Empty search should return empty results
	tags, err := service.SearchTags("")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagService_SearchTags_NoResults(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Search for something that doesn't exist
	tags, err := service.SearchTags("xyzabc123")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagService_GetAllCategories(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get all categories
	categories, err := service.GetAllCategories()

	require.NoError(t, err)
	assert.NotEmpty(t, categories)

	// Check for expected categories
	expectedCategories := []string{
		"scent_family",
		"character",
		"mood",
		"season",
		"occasion",
		"time_of_day",
		"longevity",
		"sillage",
		"ingredients",
	}

	for _, expected := range expectedCategories {
		assert.Contains(t, categories, expected,
			"Should contain category: %s", expected)
	}
}

func TestTagService_GetTagsGroupedByCategory(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get grouped tags
	grouped, err := service.GetTagsGroupedByCategory()

	require.NoError(t, err)
	assert.NotEmpty(t, grouped)

	// Verify structure
	assert.IsType(t, map[string][]string{}, grouped)

	// Check some specific categories and tags
	assert.Contains(t, grouped, "scent_family")
	assert.Contains(t, grouped["scent_family"], "citrus")
	assert.Contains(t, grouped["scent_family"], "floral")
	assert.Contains(t, grouped["scent_family"], "woody")

	assert.Contains(t, grouped, "season")
	assert.Contains(t, grouped["season"], "spring")
	assert.Contains(t, grouped["season"], "summer")
	assert.Contains(t, grouped["season"], "fall")
	assert.Contains(t, grouped["season"], "winter")

	assert.Contains(t, grouped, "character")
	assert.Contains(t, grouped["character"], "fresh")
	assert.Contains(t, grouped["character"], "warm")

	// Verify all categories have at least one tag
	for category, tags := range grouped {
		assert.NotEmpty(t, tags, "Category %s should have at least one tag", category)
	}
}

func TestTagService_GetTagsGroupedByCategory_Completeness(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get all tags
	allTags, err := service.GetAllTags()
	require.NoError(t, err)

	// Get grouped tags
	grouped, err := service.GetTagsGroupedByCategory()
	require.NoError(t, err)

	// Count total tags in grouped result
	totalGroupedTags := 0
	for _, tags := range grouped {
		totalGroupedTags += len(tags)
	}

	// Should match the total number of tags
	assert.Equal(t, len(allTags), totalGroupedTags,
		"Grouped tags should contain all tags")
}

func TestTagService_CategoryConsistency(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get all categories
	categories, err := service.GetAllCategories()
	require.NoError(t, err)

	// For each category, get tags
	for _, category := range categories {
		tags, err := service.GetTagsByCategory(category)
		require.NoError(t, err, "Should get tags for category: %s", category)

		// Verify all returned tags have the correct category
		for _, tag := range tags {
			assert.Equal(t, category, tag.Category,
				"Tag %s should belong to category %s", tag.Tag, category)
		}
	}
}

func TestTagService_SearchPartialMatch(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Search for partial matches
	testCases := []struct {
		search        string
		minResults    int
		description   string
	}{
		{"al", 1, "Should match 'floral' and other tags with 'al'"},
		{"wood", 1, "Should match 'woody'"},
		{"fresh", 1, "Should match 'fresh'"},
	}

	for _, tc := range testCases {
		t.Run(tc.search, func(t *testing.T) {
			tags, err := service.SearchTags(tc.search)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(tags), tc.minResults,
				"%s: got %d results, expected at least %d",
				tc.description, len(tags), tc.minResults)
		})
	}
}

func TestTagService_Integration_AllMethods(t *testing.T) {
	service, tdb := setupTagService(t)
	defer tdb.Teardown(t)
	defer tdb.CleanupTables(t)

	// Get all tags
	allTags, err := service.GetAllTags()
	require.NoError(t, err)
	require.NotEmpty(t, allTags)

	// Get all categories
	categories, err := service.GetAllCategories()
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	// Get grouped tags
	grouped, err := service.GetTagsGroupedByCategory()
	require.NoError(t, err)
	require.NotEmpty(t, grouped)

	// Verify consistency across all methods
	assert.Equal(t, len(categories), len(grouped),
		"Number of categories should match grouped keys")

	// For each category
	categoryTagCount := 0
	for _, category := range categories {
		// Get tags for this category
		categoryTags, err := service.GetTagsByCategory(category)
		require.NoError(t, err)

		// Should match grouped result
		assert.Equal(t, len(categoryTags), len(grouped[category]),
			"Tag count for category %s should match", category)

		categoryTagCount += len(categoryTags)
	}

	// Total tags should match
	assert.Equal(t, len(allTags), categoryTagCount,
		"Total tags across all categories should match GetAllTags")
}
