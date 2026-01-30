package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type InvitationService struct {
	repo *repository.InvitationRepository
}

func NewInvitationService(repo *repository.InvitationRepository) *InvitationService {
	return &InvitationService{repo: repo}
}

func (s *InvitationService) Create(creatorID string, email *string, expiresInDays int) (*models.Invitation, error) {
	// Generate unique code
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate invitation code: %w", err)
	}
	code := hex.EncodeToString(bytes)

	// Calculate expiration
	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	invitation := &models.Invitation{
		Code:      code,
		Email:     email,
		CreatedBy: creatorID,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := s.repo.Create(invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *InvitationService) List(creatorID string) ([]*models.Invitation, error) {
	return s.repo.ListByCreator(creatorID)
}

func (s *InvitationService) Revoke(code, creatorID string) error {
	return s.repo.Revoke(code, creatorID)
}

func (s *InvitationService) ValidateAndUse(code, email, userID string) error {
	invitation, err := s.repo.FindByCode(code)
	if err != nil {
		return fmt.Errorf("invalid invitation code")
	}

	// Check if already used
	if invitation.Used {
		return fmt.Errorf("invitation code has already been used")
	}

	// Check if expired
	if time.Now().After(invitation.ExpiresAt) {
		return fmt.Errorf("invitation code has expired")
	}

	// Check if email-specific
	if invitation.Email != nil && *invitation.Email != email {
		return fmt.Errorf("this invitation is for a different email address")
	}

	// Mark as used
	if err := s.repo.MarkAsUsed(invitation.ID, userID); err != nil {
		return err
	}

	return nil
}
