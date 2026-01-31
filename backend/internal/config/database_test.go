package config

import (
"database/sql"
"fmt"
"testing"

"github.com/yourusername/scentora-backend/internal/testutil"
)

// TestDatabaseMigrations tests that all tables are created correctly
func TestDatabaseMigrations(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)

// Test that all expected tables exist
expectedTables := []string{
"users",
"accords",
"accord_tags",
"predefined_tags",
"refresh_tokens",
"invitations",
}

for _, table := range expectedTables {
var exists bool
err := tdb.DB.QueryRow(`
SELECT EXISTS (
SELECT FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = $1
)
`, table).Scan(&exists)

if err != nil {
t.Errorf("Failed to check if table %s exists: %v", table, err)
}

if !exists {
t.Errorf("Table %s does not exist", table)
}
}
}

// TestAccordsTableStructure tests the accords table schema
func TestAccordsTableStructure(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)

// Test that volume_drops is a generated column
var result string
err := tdb.DB.QueryRow(`
SELECT column_name 
FROM information_schema.columns 
WHERE table_name = 'accords' 
AND column_name = 'volume_drops'
AND is_generated = 'ALWAYS'
`).Scan(&result)

if err == sql.ErrNoRows {
t.Error("volume_drops column is not a generated column")
} else if err != nil {
t.Errorf("Failed to check volume_drops column: %v", err)
}

// Test unique constraint on (user_id, name, pyramid_position)
var constraintName string
err = tdb.DB.QueryRow(`
SELECT constraint_name 
FROM information_schema.table_constraints 
WHERE table_name = 'accords' 
AND constraint_type = 'UNIQUE'
AND constraint_name LIKE '%name%pyramid%'
`).Scan(&constraintName)

if err == sql.ErrNoRows {
t.Error("Unique constraint on (user_id, name, pyramid_position) does not exist")
} else if err != nil {
t.Errorf("Failed to check unique constraint: %v", err)
}
}

// TestAccordsCheckConstraints tests the check constraints on accords table
func TestAccordsCheckConstraints(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create a test user
var userID string
err := tdb.DB.QueryRow(`
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3)
RETURNING id
`, "test@example.com", "testuser", "hashedpassword").Scan(&userID)

if err != nil {
t.Fatalf("Failed to create test user: %v", err)
}

// Test pyramid_position check constraint (should only allow top, middle, base)
_, err = tdb.DB.Exec(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml)
VALUES ($1, 'Test Accord', 'invalid', 10.0)
`, userID)

if err == nil {
t.Error("Expected error for invalid pyramid_position, got nil")
}

// Test volume_ml check constraint (should not allow negative values)
_, err = tdb.DB.Exec(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml)
VALUES ($1, 'Test Accord', 'top', -5.0)
`, userID)

if err == nil {
t.Error("Expected error for negative volume_ml, got nil")
}

// Test dilution_percentage check constraint (should be 0-100)
_, err = tdb.DB.Exec(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml, dilution_percentage)
VALUES ($1, 'Test Accord', 'top', 10.0, 150.0)
`, userID)

if err == nil {
t.Error("Expected error for dilution_percentage > 100, got nil")
}
}

// TestVolumeDropsCalculation tests that volume_drops is calculated correctly
func TestVolumeDropsCalculation(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create a test user
var userID string
err := tdb.DB.QueryRow(`
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3)
RETURNING id
`, "test@example.com", "testuser", "hashedpassword").Scan(&userID)

if err != nil {
t.Fatalf("Failed to create test user: %v", err)
}

testCases := []struct {
volumeMl      float64
expectedDrops int
}{
{1.0, 20},
{5.0, 100},
{10.5, 210},
{0.5, 10},
{2.75, 55},
}

for i, tc := range testCases {
var volumeDrops int
err := tdb.DB.QueryRow(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml)
VALUES ($1, $2, 'top', $3)
RETURNING volume_drops
`, userID, fmt.Sprintf("Test Accord %d", i), tc.volumeMl).Scan(&volumeDrops)

if err != nil {
t.Errorf("Failed to insert accord with volume_ml=%.2f: %v", tc.volumeMl, err)
continue
}

if volumeDrops != tc.expectedDrops {
t.Errorf("For volume_ml=%.2f, expected volume_drops=%d, got %d", tc.volumeMl, tc.expectedDrops, volumeDrops)
}
}
}

// TestAccordTagsUniqueConstraint tests the unique constraint on accord_tags
func TestAccordTagsUniqueConstraint(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create a test user
var userID string
err := tdb.DB.QueryRow(`
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3)
RETURNING id
`, "test@example.com", "testuser", "hashedpassword").Scan(&userID)

if err != nil {
t.Fatalf("Failed to create test user: %v", err)
}

// Create a test accord
var accordID string
err = tdb.DB.QueryRow(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml)
VALUES ($1, 'Test Accord', 'top', 10.0)
RETURNING id
`, userID).Scan(&accordID)

if err != nil {
t.Fatalf("Failed to create test accord: %v", err)
}

// Insert a tag
_, err = tdb.DB.Exec(`
INSERT INTO accord_tags (accord_id, tag)
VALUES ($1, 'fresh')
`, accordID)

if err != nil {
t.Fatalf("Failed to insert first tag: %v", err)
}

// Try to insert the same tag again (should fail)
_, err = tdb.DB.Exec(`
INSERT INTO accord_tags (accord_id, tag)
VALUES ($1, 'fresh')
`, accordID)

if err == nil {
t.Error("Expected error for duplicate tag, got nil")
}
}

// TestCascadeDelete tests that deleting a user cascades to accords and tags
func TestCascadeDelete(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

// Create a test user
var userID string
err := tdb.DB.QueryRow(`
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3)
RETURNING id
`, "test@example.com", "testuser", "hashedpassword").Scan(&userID)

if err != nil {
t.Fatalf("Failed to create test user: %v", err)
}

// Create a test accord
var accordID string
err = tdb.DB.QueryRow(`
INSERT INTO accords (user_id, name, pyramid_position, volume_ml)
VALUES ($1, 'Test Accord', 'top', 10.0)
RETURNING id
`, userID).Scan(&accordID)

if err != nil {
t.Fatalf("Failed to create test accord: %v", err)
}

// Insert a tag
_, err = tdb.DB.Exec(`
INSERT INTO accord_tags (accord_id, tag)
VALUES ($1, 'fresh')
`, accordID)

if err != nil {
t.Fatalf("Failed to insert tag: %v", err)
}

// Delete the user
_, err = tdb.DB.Exec(`DELETE FROM users WHERE id = $1`, userID)
if err != nil {
t.Fatalf("Failed to delete user: %v", err)
}

// Check that accord was deleted
var accordCount int
err = tdb.DB.QueryRow(`SELECT COUNT(*) FROM accords WHERE id = $1`, accordID).Scan(&accordCount)
if err != nil {
t.Fatalf("Failed to check accord count: %v", err)
}

if accordCount != 0 {
t.Error("Expected accord to be deleted via cascade, but it still exists")
}

// Check that tag was deleted
var tagCount int
err = tdb.DB.QueryRow(`SELECT COUNT(*) FROM accord_tags WHERE accord_id = $1`, accordID).Scan(&tagCount)
if err != nil {
t.Fatalf("Failed to check tag count: %v", err)
}

if tagCount != 0 {
t.Error("Expected tag to be deleted via cascade, but it still exists")
}
}
