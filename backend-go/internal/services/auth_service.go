package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/scentora-backend/internal/config"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.RefreshTokenRepository
	cfg       *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.RefreshTokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) Register(email, username, password string) (*models.AuthResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username exists
	exists, err = s.userRepo.UsernameExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Login(email, password string) (*models.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(refreshTokenString string) (*models.AuthResponse, error) {
	// Hash the token to look it up
	tokenHash := hashToken(refreshTokenString)

	// Find token in database
	token, err := s.tokenRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if expired
	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	_ = s.tokenRepo.Revoke(tokenHash)

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(refreshTokenString string) error {
	tokenHash := hashToken(refreshTokenString)
	return s.tokenRepo.Revoke(tokenHash)
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	expiresAt, err := parseDuration(s.cfg.JWTAccessExpiresIn)
	if err != nil {
		return "", err
	}

	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	expiresAt, err := parseDuration(s.cfg.JWTRefreshExpiresIn)
	if err != nil {
		return "", err
	}

	// Generate random token
	tokenString := fmt.Sprintf("%s:%d", user.ID, time.Now().UnixNano())
	tokenHash := hashToken(tokenString)

	// Store in database
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(expiresAt),
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func parseDuration(s string) (time.Duration, error) {
	// Parse duration strings like "15m", "7d", "24h"
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	value := s[:len(s)-1]
	unit := s[len(s)-1:]

	var duration time.Duration
	switch unit {
	case "m":
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		duration = d
	case "h":
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		duration = d
	case "d":
		d, err := time.ParseDuration(value + "h")
		if err != nil {
			return 0, err
		}
		duration = d * 24
	default:
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		duration = d
	}

	return duration, nil
}
