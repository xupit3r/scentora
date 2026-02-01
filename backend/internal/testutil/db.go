package testutil

import (
"database/sql"
"sync/atomic"
"fmt"
"log"
"os"
"testing"

"github.com/jmoiron/sqlx"
_ "github.com/jackc/pgx/v5/stdlib"
)

// TestDB manages a test database connection
type TestDB struct {
DB *sqlx.DB
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
// Use environment variable or default to test database
dbURL := os.Getenv("TEST_DATABASE_URL")
if dbURL == "" {
dbURL = "postgres://admin:password@localhost:5435/scentora_test?sslmode=disable"
}

db, err := sqlx.Connect("pgx", dbURL)
if err != nil {
t.Fatalf("Failed to connect to test database: %v", err)
}

// Run migrations
if err := runTestMigrations(db.DB); err != nil {
t.Fatalf("Failed to run test migrations: %v", err)
}

return &TestDB{DB: db}
}

// Teardown closes the database connection
func (tdb *TestDB) Teardown(t *testing.T) {
if err := tdb.DB.Close(); err != nil {
t.Errorf("Failed to close test database: %v", err)
}
}

// CleanupTables removes all data from tables (for test isolation)
func (tdb *TestDB) CleanupTables(t *testing.T) {
// Order matters: delete from child tables before parent tables
tables := []string{
"accord_tags",        // references accords
"accords",            // references users
"refresh_tokens",     // references users
"invitations",        // references users
"users",              // parent table
}

for _, table := range tables {
_, err := tdb.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
if err != nil {
t.Errorf("Failed to cleanup table %s: %v", table, err)
}
}
}

// runTestMigrations runs database migrations for testing
func runTestMigrations(db *sql.DB) error {
// Create users table
_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
email VARCHAR(255) UNIQUE NOT NULL,
username VARCHAR(100) NOT NULL,
password_hash VARCHAR(255) NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)
`)
if err != nil {
return fmt.Errorf("failed to create users table: %w", err)
}

// Drop old tables if they exist
_, err = db.Exec(`DROP TABLE IF EXISTS journal_entries CASCADE; DROP TABLE IF EXISTS perfumes CASCADE;`)
if err != nil {
return fmt.Errorf("failed to drop old tables: %w", err)
}

// Create accords table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS accords (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
name VARCHAR(255) NOT NULL,
pyramid_position VARCHAR(10) NOT NULL CHECK (pyramid_position IN ('top', 'middle', 'base')),
volume_ml DECIMAL(10,2) NOT NULL CHECK (volume_ml >= 0),
volume_drops INTEGER GENERATED ALWAYS AS (ROUND(volume_ml * 20)) STORED,
supplier VARCHAR(255),
purchase_date DATE,
dilution_percentage DECIMAL(5,2) CHECK (dilution_percentage >= 0 AND dilution_percentage <= 100),
notes TEXT,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE(user_id, name, pyramid_position)
)
`)
if err != nil {
return fmt.Errorf("failed to create accords table: %w", err)
}

// Create accord_tags table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS accord_tags (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE CASCADE,
tag VARCHAR(50) NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE(accord_id, tag)
)
`)
if err != nil {
return fmt.Errorf("failed to create accord_tags table: %w", err)
}

// Create predefined_tags table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS predefined_tags (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
category VARCHAR(50) NOT NULL,
tag VARCHAR(50) NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE(tag)
)
`)
if err != nil {
return fmt.Errorf("failed to create predefined_tags table: %w", err)
}

// Create refresh_tokens table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS refresh_tokens (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash VARCHAR(255) UNIQUE NOT NULL,
expires_at TIMESTAMP NOT NULL,
revoked BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
`)
if err != nil {
return fmt.Errorf("failed to create refresh_tokens table: %w", err)
}

// Create invitations table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS invitations (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
code VARCHAR(255) UNIQUE NOT NULL,
email VARCHAR(255),
created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
expires_at TIMESTAMP NOT NULL,
used BOOLEAN DEFAULT FALSE,
used_at TIMESTAMP,
used_by UUID REFERENCES users(id) ON DELETE SET NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
`)
if err != nil {
return fmt.Errorf("failed to create invitations table: %w", err)
}

// Seed predefined tags if not already seeded
var count int
err = db.QueryRow(`SELECT COUNT(*) FROM predefined_tags`).Scan(&count)
if err != nil {
return fmt.Errorf("failed to count predefined tags: %w", err)
}

if count == 0 {
err = seedPredefinedTags(db)
if err != nil {
return fmt.Errorf("failed to seed predefined tags: %w", err)
}
}

log.Println("✅ Test database migrations completed")
return nil
}

// seedPredefinedTags seeds the predefined tags for testing
func seedPredefinedTags(db *sql.DB) error {
tags := []struct {
category string
tag      string
}{
// Scent families
{"scent_family", "citrus"},
{"scent_family", "floral"},
{"scent_family", "woody"},
{"scent_family", "oriental"},
{"scent_family", "fruity"},
{"scent_family", "aquatic"},
{"scent_family", "gourmand"},
{"scent_family", "chypre"},
{"scent_family", "fougere"},

// Character
{"character", "fresh"},
{"character", "warm"},
{"character", "sweet"},
{"character", "spicy"},
{"character", "earthy"},
{"character", "powdery"},
{"character", "clean"},
{"character", "rich"},
{"character", "smoky"},

// Mood
{"mood", "uplifting"},
{"mood", "calming"},
{"mood", "sensual"},
{"mood", "energizing"},
{"mood", "mysterious"},
{"mood", "comforting"},
{"mood", "elegant"},
{"mood", "playful"},

// Season
{"season", "spring"},
{"season", "summer"},
{"season", "fall"},
{"season", "winter"},
{"season", "all-season"},

// Occasion
{"occasion", "casual"},
{"occasion", "formal"},
{"occasion", "office"},
{"occasion", "evening"},
{"occasion", "date-night"},
{"occasion", "special-occasion"},

// Time of day
{"time_of_day", "morning"},
{"time_of_day", "afternoon"},
{"time_of_day", "evening"},
{"time_of_day", "night"},

// Longevity
{"longevity", "very-weak"},
{"longevity", "weak"},
{"longevity", "moderate"},
{"longevity", "long-lasting"},
{"longevity", "eternal"},

// Sillage (projection)
{"sillage", "intimate"},
{"sillage", "moderate"},
{"sillage", "strong"},
{"sillage", "enormous"},

// Common ingredients
{"ingredients", "bergamot"},
{"ingredients", "lavender"},
{"ingredients", "sandalwood"},
{"ingredients", "vanilla"},
{"ingredients", "rose"},
{"ingredients", "jasmine"},
{"ingredients", "patchouli"},
{"ingredients", "musk"},
{"ingredients", "amber"},
}

for _, t := range tags {
_, err := db.Exec(`
INSERT INTO predefined_tags (category, tag)
VALUES ($1, $2)
ON CONFLICT (tag) DO NOTHING
`, t.category, t.tag)
if err != nil {
return fmt.Errorf("failed to insert tag %s: %w", t.tag, err)
}
}

return nil
}

var testCounter uint64

// UniqueEmail generates a unique email for testing
func UniqueEmail(prefix string) string {
count := atomic.AddUint64(&testCounter, 1)
return fmt.Sprintf("%s_%d@test.example.com", prefix, count)
}
