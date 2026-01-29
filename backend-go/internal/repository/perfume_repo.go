package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yourusername/scentora-backend/internal/models"
)

type PerfumeRepository struct {
	db *sqlx.DB
}

func NewPerfumeRepository(db *sqlx.DB) *PerfumeRepository {
	return &PerfumeRepository{db: db}
}

func (r *PerfumeRepository) Create(perfume *models.Perfume) error {
	query := `
		INSERT INTO perfumes (user_id, name, designer, year, concentration, top_notes, middle_notes, base_notes, description, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		perfume.UserID,
		perfume.Name,
		perfume.Designer,
		perfume.Year,
		perfume.Concentration,
		pq.Array(perfume.TopNotes),
		pq.Array(perfume.MiddleNotes),
		pq.Array(perfume.BaseNotes),
		perfume.Description,
		perfume.ImageURL,
	).Scan(&perfume.ID, &perfume.CreatedAt, &perfume.UpdatedAt)
}

func (r *PerfumeRepository) FindByID(id, userID string) (*models.Perfume, error) {
	var perfume models.Perfume
	query := `
		SELECT id, user_id, name, designer, year, concentration, top_notes, middle_notes, base_notes, description, image_url, created_at, updated_at
		FROM perfumes
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&perfume, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &perfume, nil
}

func (r *PerfumeRepository) List(userID string, filters map[string]string) ([]*models.Perfume, error) {
	query := `
		SELECT id, user_id, name, designer, year, concentration, top_notes, middle_notes, base_notes, description, image_url, created_at, updated_at
		FROM perfumes
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if search := filters["search"]; search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR designer ILIKE $%d OR description ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+search+"%")
	}

	if year := filters["year"]; year != "" {
		argCount++
		query += fmt.Sprintf(" AND year = $%d", argCount)
		args = append(args, year)
	}

	if concentration := filters["concentration"]; concentration != "" {
		argCount++
		query += fmt.Sprintf(" AND concentration = $%d", argCount)
		args = append(args, concentration)
	}

	if note := filters["note"]; note != "" {
		argCount++
		query += fmt.Sprintf(" AND ($%d = ANY(top_notes) OR $%d = ANY(middle_notes) OR $%d = ANY(base_notes))", argCount, argCount, argCount)
		args = append(args, note)
	}

	query += " ORDER BY created_at DESC"

	var perfumes []*models.Perfume
	err := r.db.Select(&perfumes, query, args...)
	if err != nil {
		return nil, err
	}
	return perfumes, nil
}

func (r *PerfumeRepository) Update(id, userID string, updates map[string]interface{}) (*models.Perfume, error) {
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
		"UPDATE perfumes SET %s WHERE id = $%d AND user_id = $%d RETURNING id, user_id, name, designer, year, concentration, top_notes, middle_notes, base_notes, description, image_url, created_at, updated_at",
		strings.Join(setParts, ", "),
		whereID,
		whereUserID,
	)

	args = append(args, id, userID)

	var perfume models.Perfume
	err := r.db.QueryRowx(query, args...).StructScan(&perfume)
	if err != nil {
		return nil, err
	}
	return &perfume, nil
}

func (r *PerfumeRepository) Delete(id, userID string) error {
	query := `DELETE FROM perfumes WHERE id = $1 AND user_id = $2`
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

func (r *PerfumeRepository) GetAllNotes(userID string) ([]string, error) {
	query := `
		SELECT DISTINCT unnest(top_notes || middle_notes || base_notes) as note
		FROM perfumes
		WHERE user_id = $1
		ORDER BY note
	`
	var notes []string
	err := r.db.Select(&notes, query, userID)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *PerfumeRepository) Count(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM perfumes WHERE user_id = $1`
	err := r.db.Get(&count, query, userID)
	return count, err
}
