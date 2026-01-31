package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yourusername/scentora-backend/internal/models"
)

type AccordRepository struct {
	db *sqlx.DB
}

func NewAccordRepository(db *sqlx.DB) *AccordRepository {
	return &AccordRepository{db: db}
}

// Create creates a new accord
func (r *AccordRepository) Create(accord *models.Accord) error {
	query := `
		INSERT INTO accords (
			user_id, name, pyramid_position, volume_ml, 
			supplier, purchase_date, dilution_percentage, notes,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, volume_drops, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		accord.UserID,
		accord.Name,
		accord.PyramidPosition,
		accord.VolumeMl,
		accord.Supplier,
		accord.PurchaseDate,
		accord.DilutionPercentage,
		accord.Notes,
	).Scan(&accord.ID, &accord.VolumeDrops, &accord.CreatedAt, &accord.UpdatedAt)
}

// FindByID retrieves an accord by ID
func (r *AccordRepository) FindByID(id, userID string) (*models.Accord, error) {
	var accord models.Accord
	query := `
		SELECT id, user_id, name, pyramid_position, volume_ml, volume_drops,
		       supplier, purchase_date, dilution_percentage, notes,
		       created_at, updated_at
		FROM accords
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.Get(&accord, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &accord, nil
}

// List retrieves all accords for a user
func (r *AccordRepository) List(userID string) ([]*models.Accord, error) {
	var accords []*models.Accord
	query := `
		SELECT id, user_id, name, pyramid_position, volume_ml, volume_drops,
		       supplier, purchase_date, dilution_percentage, notes,
		       created_at, updated_at
		FROM accords
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&accords, query, userID)
	if err != nil {
		return nil, err
	}
	return accords, nil
}

// Update updates an accord
func (r *AccordRepository) Update(accord *models.Accord) error {
	query := `
		UPDATE accords
		SET name = $1, pyramid_position = $2, volume_ml = $3,
		    supplier = $4, purchase_date = $5, dilution_percentage = $6,
		    notes = $7, updated_at = NOW()
		WHERE id = $8 AND user_id = $9
		RETURNING volume_drops, updated_at
	`
	result := r.db.QueryRow(
		query,
		accord.Name,
		accord.PyramidPosition,
		accord.VolumeMl,
		accord.Supplier,
		accord.PurchaseDate,
		accord.DilutionPercentage,
		accord.Notes,
		accord.ID,
		accord.UserID,
	)

	err := result.Scan(&accord.VolumeDrops, &accord.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an accord
func (r *AccordRepository) Delete(id, userID string) error {
	query := `DELETE FROM accords WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("accord not found")
	}
	return nil
}

// Filter retrieves accords with optional filters
func (r *AccordRepository) Filter(userID string, position *string, minVolume, maxVolume *float64, supplier, search *string) ([]*models.Accord, error) {
	var accords []*models.Accord

	query := `
		SELECT id, user_id, name, pyramid_position, volume_ml, volume_drops,
		       supplier, purchase_date, dilution_percentage, notes,
		       created_at, updated_at
		FROM accords
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	if position != nil {
		query += fmt.Sprintf(" AND pyramid_position = $%d", argIndex)
		args = append(args, *position)
		argIndex++
	}

	if minVolume != nil {
		query += fmt.Sprintf(" AND volume_ml >= $%d", argIndex)
		args = append(args, *minVolume)
		argIndex++
	}

	if maxVolume != nil {
		query += fmt.Sprintf(" AND volume_ml <= $%d", argIndex)
		args = append(args, *maxVolume)
		argIndex++
	}

	if supplier != nil {
		query += fmt.Sprintf(" AND supplier ILIKE $%d", argIndex)
		args = append(args, "%"+*supplier+"%")
		argIndex++
	}

	if search != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+*search+"%")
		argIndex++
	}

	query += " ORDER BY name"

	err := r.db.Select(&accords, query, args...)
	if err != nil {
		return nil, err
	}

	return accords, nil
}


// ListWithFilters retrieves accords with optional filters
func (r *AccordRepository) ListWithFilters(userID string, filters map[string]interface{}) ([]*models.Accord, error) {
	var accords []*models.Accord

	query := `
		SELECT id, user_id, name, pyramid_position, volume_ml, volume_drops,
		       supplier, purchase_date, dilution_percentage, notes,
		       created_at, updated_at
		FROM accords
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	// Apply filters
	if position, ok := filters["position"].(string); ok && position != "" {
		query += fmt.Sprintf(" AND pyramid_position = $%d", argIndex)
		args = append(args, position)
		argIndex++
	}

	if minVolume, ok := filters["min_volume"].(float64); ok {
		query += fmt.Sprintf(" AND volume_ml >= $%d", argIndex)
		args = append(args, minVolume)
		argIndex++
	}

	if maxVolume, ok := filters["max_volume"].(float64); ok {
		query += fmt.Sprintf(" AND volume_ml <= $%d", argIndex)
		args = append(args, maxVolume)
		argIndex++
	}

	if supplier, ok := filters["supplier"].(string); ok && supplier != "" {
		query += fmt.Sprintf(" AND supplier ILIKE $%d", argIndex)
		args = append(args, "%"+supplier+"%")
		argIndex++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR notes ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	err := r.db.Select(&accords, query, args...)
	if err != nil {
		return nil, err
	}
	return accords, nil
}

// GetTagsForAccord retrieves all tags for an accord
func (r *AccordRepository) GetTagsForAccord(accordID string) ([]string, error) {
	var tags []string
	query := `
		SELECT tag
		FROM accord_tags
		WHERE accord_id = $1
		ORDER BY tag ASC
	`
	err := r.db.Select(&tags, query, accordID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// AddTag adds a tag to an accord
func (r *AccordRepository) AddTag(accordID, tag string) error {
	query := `
		INSERT INTO accord_tags (accord_id, tag, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := r.db.Exec(query, accordID, strings.ToLower(strings.TrimSpace(tag)))
	return err
}

// RemoveTag removes a tag from an accord
func (r *AccordRepository) RemoveTag(accordID, tag string) error {
	query := `DELETE FROM accord_tags WHERE accord_id = $1 AND tag = $2`
	result, err := r.db.Exec(query, accordID, strings.ToLower(strings.TrimSpace(tag)))
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}
	return nil
}

// SetTags replaces all tags for an accord
func (r *AccordRepository) SetTags(accordID string, tags []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing tags
	_, err = tx.Exec(`DELETE FROM accord_tags WHERE accord_id = $1`, accordID)
	if err != nil {
		return err
	}

	// Insert new tags
	if len(tags) > 0 {
		query := `INSERT INTO accord_tags (accord_id, tag, created_at) VALUES ($1, $2, NOW())`
		for _, tag := range tags {
			cleanTag := strings.ToLower(strings.TrimSpace(tag))
			if cleanTag != "" {
				_, err = tx.Exec(query, accordID, cleanTag)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

// GetAccordsByTags retrieves accords that have ALL specified tags
func (r *AccordRepository) GetAccordsByTags(userID string, tags []string) ([]*models.Accord, error) {
	if len(tags) == 0 {
		return r.List(userID)
	}

	var accords []*models.Accord
	
	// Build query to find accords with all specified tags
	query := `
		SELECT DISTINCT a.id, a.user_id, a.name, a.pyramid_position, a.volume_ml, a.volume_drops,
		       a.supplier, a.purchase_date, a.dilution_percentage, a.notes,
		       a.created_at, a.updated_at
		FROM accords a
		INNER JOIN accord_tags at ON a.id = at.accord_id
		WHERE a.user_id = $1 AND at.tag = ANY($2)
		GROUP BY a.id
		HAVING COUNT(DISTINCT at.tag) = $3
		ORDER BY a.created_at DESC
	`

	err := r.db.Select(&accords, query, userID, tags, len(tags))
	if err != nil {
		return nil, err
	}
	return accords, nil
}
