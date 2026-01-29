package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port    string
	NodeEnv string

	// Database
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// JWT
	JWTSecret           string
	JWTAccessExpiresIn  string
	JWTRefreshExpiresIn string

	// CORS
	CORSAllowedOrigins []string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Port:                getEnv("PORT", "3000"),
		NodeEnv:             getEnv("NODE_ENV", "development"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "admin"),
		DBPassword:          getEnv("DB_PASSWORD", "password"),
		DBName:              getEnv("DB_NAME", "scentora"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		JWTAccessExpiresIn:  getEnv("JWT_ACCESS_EXPIRES_IN", "15m"),
		JWTRefreshExpiresIn: getEnv("JWT_REFRESH_EXPIRES_IN", "7d"),
		RateLimitRequests:   100,
		RateLimitWindow:     getEnv("RATE_LIMIT_WINDOW", "1m"),
	}

	// Parse CORS origins
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	cfg.CORSAllowedOrigins = strings.Split(originsStr, ",")

	// Build DATABASE_URL if not provided
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
		)
	}

	// Validate required fields
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
