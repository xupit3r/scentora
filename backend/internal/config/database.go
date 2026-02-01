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

	// ========== PHASE 10: RECIPE SYSTEM MIGRATIONS ==========
	
	// Add validate_recipe_volumes to users table
	_, err = db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS validate_recipe_volumes BOOLEAN DEFAULT FALSE;
	`)
	if err != nil {
		return fmt.Errorf("failed to alter users table: %w", err)
	}

	// Create recipes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			target_volume_ml DECIMAL(10,2) NOT NULL CHECK (target_volume_ml > 0),
			status VARCHAR(20) NOT NULL DEFAULT 'draft' 
				CHECK (status IN ('draft', 'in_progress', 'tested', 'finalized', 'archived')),
			active_version_id UUID,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, name)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipes table: %w", err)
	}

	// Create recipe_versions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_versions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			version_number INTEGER NOT NULL,
			name VARCHAR(100) NOT NULL,
			notes TEXT,
			is_active BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(recipe_id, version_number)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_versions table: %w", err)
	}

	// Add foreign key constraint for active_version_id (after both tables exist)
	_, err = db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_recipes_active_version'
			) THEN
				ALTER TABLE recipes 
					ADD CONSTRAINT fk_recipes_active_version 
					FOREIGN KEY (active_version_id) 
					REFERENCES recipe_versions(id) ON DELETE SET NULL;
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("failed to add active_version foreign key: %w", err)
	}

	// Create recipe_ingredients table
	// CRITICAL: ON DELETE RESTRICT prevents accord deletion if used in recipe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_ingredients (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			version_id UUID NOT NULL REFERENCES recipe_versions(id) ON DELETE CASCADE,
			accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE RESTRICT,
			quantity_ml DECIMAL(10,2) NOT NULL CHECK (quantity_ml > 0),
			quantity_drops INTEGER GENERATED ALWAYS AS (ROUND(quantity_ml * 20)) STORED,
			percentage DECIMAL(5,2) CHECK (percentage >= 0 AND percentage <= 100),
			notes TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(version_id, accord_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_ingredients table: %w", err)
	}

	// Create recipe_notes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_notes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			version_id UUID REFERENCES recipe_versions(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			note_type VARCHAR(20) DEFAULT 'general' 
				CHECK (note_type IN ('general', 'testing', 'observation')),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_notes table: %w", err)
	}

	// Create recipe_tags table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_tags (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			tag VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(recipe_id, tag)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_tags table: %w", err)
	}

	// Create recipe_collections table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_collections (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, name)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_collections table: %w", err)
	}

	// Create recipe_collection_members table (join table)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recipe_collection_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			collection_id UUID NOT NULL REFERENCES recipe_collections(id) ON DELETE CASCADE,
			recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			added_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(collection_id, recipe_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe_collection_members table: %w", err)
	}

	// Create indexes for recipe tables
	_, err = db.Exec(`
		-- Recipe indexes
		CREATE INDEX IF NOT EXISTS idx_recipes_user_id ON recipes(user_id);
		CREATE INDEX IF NOT EXISTS idx_recipes_status ON recipes(status);
		CREATE INDEX IF NOT EXISTS idx_recipes_created_at ON recipes(created_at DESC);
		
		-- Recipe version indexes
		CREATE INDEX IF NOT EXISTS idx_recipe_versions_recipe_id ON recipe_versions(recipe_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_versions_is_active ON recipe_versions(is_active);
		
		-- Recipe ingredient indexes
		CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_version_id ON recipe_ingredients(version_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_accord_id ON recipe_ingredients(accord_id);
		
		-- Recipe note indexes
		CREATE INDEX IF NOT EXISTS idx_recipe_notes_recipe_id ON recipe_notes(recipe_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_notes_version_id ON recipe_notes(version_id);
		
		-- Recipe tag indexes
		CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe_id ON recipe_tags(recipe_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_tags_tag ON recipe_tags(tag);
		
		-- Recipe collection indexes
		CREATE INDEX IF NOT EXISTS idx_recipe_collections_user_id ON recipe_collections(user_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_collection_members_collection_id ON recipe_collection_members(collection_id);
		CREATE INDEX IF NOT EXISTS idx_recipe_collection_members_recipe_id ON recipe_collection_members(recipe_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create recipe indexes: %w", err)
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
