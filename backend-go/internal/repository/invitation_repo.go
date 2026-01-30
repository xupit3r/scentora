package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type InvitationRepository struct {
	db *sqlx.DB
}

func NewInvitationRepository(db *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(invitation *models.Invitation) error {
	query := `
		INSERT INTO invitations (code, email, created_by, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		invitation.Code,
		invitation.Email,
		invitation.CreatedBy,
		invitation.ExpiresAt,
		invitation.Used,
	).Scan(&invitation.ID, &invitation.CreatedAt)
}

func (r *InvitationRepository) FindByCode(code string) (*models.Invitation, error) {
	var invitation models.Invitation
	query := `
		SELECT id, code, email, created_by, expires_at, used, used_at, used_by, created_at
		FROM invitations
		WHERE code = $1
	`
	err := r.db.Get(&invitation, query, code)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *InvitationRepository) ListByCreator(creatorID string) ([]*models.Invitation, error) {
	var invitations []*models.Invitation
	query := `
		SELECT id, code, email, created_by, expires_at, used, used_at, used_by, created_at
		FROM invitations
		WHERE created_by = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&invitations, query, creatorID)
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *InvitationRepository) MarkAsUsed(id, usedBy string) error {
	query := `
		UPDATE invitations
		SET used = TRUE, used_at = NOW(), used_by = $2
		WHERE id = $1
	`
	result, err := r.db.Exec(query, id, usedBy)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *InvitationRepository) Revoke(code, creatorID string) error {
	query := `
		UPDATE invitations
		SET used = TRUE, used_at = NOW()
		WHERE code = $1 AND created_by = $2
	`
	result, err := r.db.Exec(query, code, creatorID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("invitation not found or not authorized")
	}
	return nil
}
