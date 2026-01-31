package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type PredefinedTagRepository struct {
	db *sqlx.DB
}

func NewPredefinedTagRepository(db *sqlx.DB) *PredefinedTagRepository {
	return &PredefinedTagRepository{db: db}
}

// GetAll retrieves all predefined tags
func (r *PredefinedTagRepository) GetAll() ([]*models.PredefinedTag, error) {
	var tags []*models.PredefinedTag
	query := `
		SELECT id, category, tag, created_at
		FROM predefined_tags
		ORDER BY category, tag
	`
	err := r.db.Select(&tags, query)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetByCategory retrieves predefined tags by category
func (r *PredefinedTagRepository) GetByCategory(category string) ([]*models.PredefinedTag, error) {
	var tags []*models.PredefinedTag
	query := `
		SELECT id, category, tag, created_at
		FROM predefined_tags
		WHERE category = $1
		ORDER BY tag
	`
	err := r.db.Select(&tags, query, category)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Search searches for tags by partial match
func (r *PredefinedTagRepository) Search(search string) ([]*models.PredefinedTag, error) {
	var tags []*models.PredefinedTag
	query := `
		SELECT id, category, tag, created_at
		FROM predefined_tags
		WHERE tag ILIKE $1
		ORDER BY tag
		LIMIT 20
	`
	err := r.db.Select(&tags, query, search+"%")
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAllCategories retrieves all unique categories
func (r *PredefinedTagRepository) GetAllCategories() ([]string, error) {
	var categories []string
	query := `
		SELECT DISTINCT category
		FROM predefined_tags
		ORDER BY category
	`
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
