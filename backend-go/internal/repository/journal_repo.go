package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type JournalRepository struct {
	db *sqlx.DB
}

func NewJournalRepository(db *sqlx.DB) *JournalRepository {
	return &JournalRepository{db: db}
}

func (r *JournalRepository) Create(entry *models.JournalEntry) error {
	query := `
		INSERT INTO journal_entries (user_id, perfume_id, date, content, rating, occasion, weather, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		entry.UserID,
		entry.PerfumeID,
		entry.Date,
		entry.Content,
		entry.Rating,
		entry.Occasion,
		entry.Weather,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

func (r *JournalRepository) FindByID(id, userID string) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	query := `
		SELECT id, user_id, perfume_id, date, content, rating, occasion, weather, created_at, updated_at
		FROM journal_entries
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&entry, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *JournalRepository) ListByPerfume(perfumeID, userID string) ([]*models.JournalEntry, error) {
	var entries []*models.JournalEntry
	query := `
		SELECT id, user_id, perfume_id, date, content, rating, occasion, weather, created_at, updated_at
		FROM journal_entries
		WHERE perfume_id = $1 AND user_id = $2
		ORDER BY date DESC, created_at DESC
	`
	err := r.db.Select(&entries, query, perfumeID, userID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *JournalRepository) Update(id, userID string, updates map[string]interface{}) (*models.JournalEntry, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argCount := 0

	for key, value := range updates {
		argCount++
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argCount))
		args = append(args, value)
	}

	argCount++
	whereID := argCount
	argCount++
	whereUserID := argCount

	query := fmt.Sprintf(
		"UPDATE journal_entries SET %s WHERE id = $%d AND user_id = $%d RETURNING id, user_id, perfume_id, date, content, rating, occasion, weather, created_at, updated_at",
		strings.Join(setParts, ", "),
		whereID,
		whereUserID,
	)

	args = append(args, id, userID)

	var entry models.JournalEntry
	err := r.db.QueryRowx(query, args...).StructScan(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *JournalRepository) Delete(id, userID string) error {
	query := `DELETE FROM journal_entries WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *JournalRepository) Count(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM journal_entries WHERE user_id = $1`
	err := r.db.Get(&count, query, userID)
	return count, err
}

func (r *JournalRepository) ListAll(userID string) ([]*models.JournalEntry, error) {
	var entries []*models.JournalEntry
	query := `
		SELECT id, user_id, perfume_id, date, content, rating, occasion, weather, created_at, updated_at
		FROM journal_entries
		WHERE user_id = $1
		ORDER BY date DESC
	`
	err := r.db.Select(&entries, query, userID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
