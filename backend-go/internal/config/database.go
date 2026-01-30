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

	// Create perfumes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS perfumes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			designer VARCHAR(255) NOT NULL,
			year INTEGER,
			concentration VARCHAR(50),
			top_notes TEXT[],
			middle_notes TEXT[],
			base_notes TEXT[],
			description TEXT,
			image_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create perfumes table: %w", err)
	}

	// Create indexes for perfumes
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_perfumes_user_id ON perfumes(user_id);
		CREATE INDEX IF NOT EXISTS idx_perfumes_created_at ON perfumes(created_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create perfumes indexes: %w", err)
	}

	// Create journal_entries table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS journal_entries (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			perfume_id UUID NOT NULL REFERENCES perfumes(id) ON DELETE CASCADE,
			date DATE NOT NULL,
			content TEXT NOT NULL,
			rating INTEGER CHECK (rating >= 1 AND rating <= 10),
			occasion VARCHAR(100),
			weather VARCHAR(100),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create journal_entries table: %w", err)
	}

	// Create indexes for journal entries
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_journal_user_perfume ON journal_entries(user_id, perfume_id);
		CREATE INDEX IF NOT EXISTS idx_journal_date ON journal_entries(date DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create journal_entries indexes: %w", err)
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
