package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDatabase(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")

	// Run migrations
	if err := runMigrations(db.DB); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
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

	// Drop old perfume and journal tables (Phase 8 migration)
	_, err = db.Exec(`
		DROP TABLE IF EXISTS journal_entries CASCADE;
		DROP TABLE IF EXISTS perfumes CASCADE;
	`)
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

	// Create indexes for accords
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_accords_user_id ON accords(user_id);
		CREATE INDEX IF NOT EXISTS idx_accords_position ON accords(pyramid_position);
		CREATE INDEX IF NOT EXISTS idx_accords_created_at ON accords(created_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create accords indexes: %w", err)
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

	// Create indexes for accord_tags
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_accord_tags_accord_id ON accord_tags(accord_id);
		CREATE INDEX IF NOT EXISTS idx_accord_tags_tag ON accord_tags(tag);
	`)
	if err != nil {
		return fmt.Errorf("failed to create accord_tags indexes: %w", err)
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

	// Create index for predefined_tags
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_predefined_tags_category ON predefined_tags(category);
	`)
	if err != nil {
		return fmt.Errorf("failed to create predefined_tags indexes: %w", err)
	}

	// Seed predefined tags
	if err := seedPredefinedTags(db); err != nil {
		return fmt.Errorf("failed to seed predefined tags: %w", err)
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

	// Create indexes for refresh tokens
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
	`)
	if err != nil {
		return fmt.Errorf("failed to create refresh_tokens indexes: %w", err)
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

	// Create indexes for invitations
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_invitations_code ON invitations(code);
		CREATE INDEX IF NOT EXISTS idx_invitations_created_by ON invitations(created_by);
	`)
	if err != nil {
		return fmt.Errorf("failed to create invitations indexes: %w", err)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

func seedPredefinedTags(db *sql.DB) error {
	// Check if tags already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM predefined_tags").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check predefined tags: %w", err)
	}

	// Skip if tags already seeded
	if count > 0 {
		log.Printf("✅ Predefined tags already seeded (%d tags)", count)
		return nil
	}

	// Character tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('character', 'fresh'),
			('character', 'warm'),
			('character', 'cool'),
			('character', 'dry'),
			('character', 'powdery'),
			('character', 'creamy'),
			('character', 'sharp'),
			('character', 'soft'),
			('character', 'rich'),
			('character', 'light')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed character tags: %w", err)
	}

	// Mood tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('mood', 'romantic'),
			('mood', 'sensual'),
			('mood', 'energetic'),
			('mood', 'calming'),
			('mood', 'mysterious'),
			('mood', 'playful'),
			('mood', 'sophisticated'),
			('mood', 'innocent')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed mood tags: %w", err)
	}

	// Season tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('season', 'spring'),
			('season', 'summer'),
			('season', 'autumn'),
			('season', 'winter')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed season tags: %w", err)
	}

	// Time tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('time', 'morning'),
			('time', 'afternoon'),
			('time', 'evening'),
			('time', 'night')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed time tags: %w", err)
	}

	// Intensity tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('intensity', 'subtle'),
			('intensity', 'moderate'),
			('intensity', 'strong'),
			('intensity', 'intense'),
			('intensity', 'bold')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed intensity tags: %w", err)
	}

	// Quality tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('quality', 'clean'),
			('quality', 'dirty'),
			('quality', 'animalic'),
			('quality', 'synthetic'),
			('quality', 'natural'),
			('quality', 'modern'),
			('quality', 'vintage')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed quality tags: %w", err)
	}

	// Scent family tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('scent_family', 'floral'),
			('scent_family', 'fruity'),
			('scent_family', 'woody'),
			('scent_family', 'oriental'),
			('scent_family', 'citrus'),
			('scent_family', 'aromatic'),
			('scent_family', 'spicy'),
			('scent_family', 'gourmand')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed scent family tags: %w", err)
	}

	// Texture tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('texture', 'smooth'),
			('texture', 'rough'),
			('texture', 'silky'),
			('texture', 'velvety'),
			('texture', 'airy'),
			('texture', 'dense')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed texture tags: %w", err)
	}

	// Style tags
	_, err = db.Exec(`
		INSERT INTO predefined_tags (category, tag) VALUES
			('style', 'casual'),
			('style', 'formal'),
			('style', 'sporty'),
			('style', 'elegant'),
			('style', 'edgy')
	`)
	if err != nil {
		return fmt.Errorf("failed to seed style tags: %w", err)
	}

	log.Println("✅ Predefined tags seeded successfully (57 tags)")
	return nil
}
