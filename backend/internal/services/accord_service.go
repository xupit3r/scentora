package services

import (
	"errors"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type AccordService struct {
	accordRepo *repository.AccordRepository
	tagRepo    *repository.PredefinedTagRepository
}

func NewAccordService(accordRepo *repository.AccordRepository, tagRepo *repository.PredefinedTagRepository) *AccordService {
	return &AccordService{
		accordRepo: accordRepo,
		tagRepo:    tagRepo,
	}
}

// CreateAccord creates a new accord with validation
func (s *AccordService) CreateAccord(userID string, req *models.CreateAccordRequest) (*models.Accord, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("accord name is required")
	}
	if req.PyramidPosition == "" {
		return nil, errors.New("pyramid position is required")
	}
	if req.PyramidPosition != "top" && req.PyramidPosition != "middle" && req.PyramidPosition != "base" {
		return nil, errors.New("pyramid position must be one of: top, middle, base")
	}
	if req.VolumeMl <= 0 {
		return nil, errors.New("volume must be greater than 0")
	}
	if req.DilutionPercentage != nil && (*req.DilutionPercentage < 0 || *req.DilutionPercentage > 100) {
		return nil, errors.New("dilution percentage must be between 0 and 100")
	}

	accord := &models.Accord{
		UserID:             userID,
		Name:               req.Name,
		PyramidPosition:    req.PyramidPosition,
		VolumeMl:           req.VolumeMl,
		Supplier:           req.Supplier,
		PurchaseDate:       req.PurchaseDate,
		DilutionPercentage: req.DilutionPercentage,
		Notes:              req.Notes,
	}

	err := s.accordRepo.Create(accord)
	if err != nil {
		return nil, fmt.Errorf("failed to create accord: %w", err)
	}

	// Add initial tags if provided
	if len(req.Tags) > 0 {
		err = s.accordRepo.SetTags(accord.ID, req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to set tags: %w", err)
		}
	}

	return accord, nil
}

// GetAccord retrieves an accord by ID
func (s *AccordService) GetAccord(accordID, userID string) (*models.Accord, error) {
	accord, err := s.accordRepo.FindByID(accordID, userID)
	if err != nil {
		return nil, fmt.Errorf("accord not found: %w", err)
	}

	// Load tags
	tags, err := s.accordRepo.GetTagsForAccord(accordID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	accord.Tags = tags

	return accord, nil
}

// ListAccords retrieves all accords for a user
func (s *AccordService) ListAccords(userID string) ([]*models.Accord, error) {
	accords, err := s.accordRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list accords: %w", err)
	}

	// Load tags for each accord
	for _, accord := range accords {
		tags, err := s.accordRepo.GetTagsForAccord(accord.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags for accord %s: %w", accord.ID, err)
		}
		accord.Tags = tags
	}

	return accords, nil
}

// UpdateAccord updates an existing accord
func (s *AccordService) UpdateAccord(accordID, userID string, req *models.UpdateAccordRequest) (*models.Accord, error) {
	// Verify accord exists and belongs to user
	accord, err := s.accordRepo.FindByID(accordID, userID)
	if err != nil {
		return nil, fmt.Errorf("accord not found: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		accord.Name = *req.Name
	}
	if req.PyramidPosition != nil {
		if *req.PyramidPosition != "top" && *req.PyramidPosition != "middle" && *req.PyramidPosition != "base" {
			return nil, errors.New("pyramid position must be one of: top, middle, base")
		}
		accord.PyramidPosition = *req.PyramidPosition
	}
	if req.VolumeMl != nil {
		if *req.VolumeMl <= 0 {
			return nil, errors.New("volume must be greater than 0")
		}
		accord.VolumeMl = *req.VolumeMl
	}
	if req.Supplier != nil {
		accord.Supplier = req.Supplier
	}
	if req.PurchaseDate != nil {
		accord.PurchaseDate = req.PurchaseDate
	}
	if req.DilutionPercentage != nil {
		if *req.DilutionPercentage < 0 || *req.DilutionPercentage > 100 {
			return nil, errors.New("dilution percentage must be between 0 and 100")
		}
		accord.DilutionPercentage = req.DilutionPercentage
	}
	if req.Notes != nil {
		accord.Notes = req.Notes
	}

	err = s.accordRepo.Update(accord)
	if err != nil {
		return nil, fmt.Errorf("failed to update accord: %w", err)
	}

	// Update tags if provided
	if req.Tags != nil {
		err = s.accordRepo.SetTags(accord.ID, *req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
	}

	// Reload to get updated data including tags
	return s.GetAccord(accordID, userID)
}

// DeleteAccord deletes an accord
func (s *AccordService) DeleteAccord(accordID, userID string) error {
	// Verify accord exists and belongs to user
	_, err := s.accordRepo.FindByID(accordID, userID)
	if err != nil {
		return fmt.Errorf("accord not found: %w", err)
	}

	err = s.accordRepo.Delete(accordID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete accord: %w", err)
	}

	return nil
}

// AddTagToAccord adds a single tag to an accord
func (s *AccordService) AddTagToAccord(accordID, userID, tag string) error {
	// Verify accord exists and belongs to user
	_, err := s.accordRepo.FindByID(accordID, userID)
	if err != nil {
		return fmt.Errorf("accord not found: %w", err)
	}

	err = s.accordRepo.AddTag(accordID, tag)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

// RemoveTagFromAccord removes a single tag from an accord
func (s *AccordService) RemoveTagFromAccord(accordID, userID, tag string) error {
	// Verify accord exists and belongs to user
	_, err := s.accordRepo.FindByID(accordID, userID)
	if err != nil {
		return fmt.Errorf("accord not found: %w", err)
	}

	err = s.accordRepo.RemoveTag(accordID, tag)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

// SearchAccords searches accords with filters
func (s *AccordService) SearchAccords(userID string, position *string, minVolume, maxVolume *float64, supplier, search *string, tags []string) ([]*models.Accord, error) {
	var accords []*models.Accord
	var err error

	// If tags are provided, search by tags
	if len(tags) > 0 {
		accords, err = s.accordRepo.GetAccordsByTags(userID, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to search accords by tags: %w", err)
		}
	} else {
		// Otherwise use the filter method
		accords, err = s.accordRepo.Filter(userID, position, minVolume, maxVolume, supplier, search)
		if err != nil {
			return nil, fmt.Errorf("failed to filter accords: %w", err)
		}
	}

	// Load tags for each accord
	for _, accord := range accords {
		accordTags, err := s.accordRepo.GetTagsForAccord(accord.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags for accord %s: %w", accord.ID, err)
		}
		accord.Tags = accordTags
	}

	return accords, nil
}

// GetStatistics returns statistics about the user's accord collection
func (s *AccordService) GetStatistics(userID string) (*models.AccordStatistics, error) {
	// Get all accords for the user
	accords, err := s.accordRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accords: %w", err)
	}

	// Load tags for each accord
	for _, accord := range accords {
		tags, err := s.accordRepo.GetTagsForAccord(accord.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags: %w", err)
		}
		accord.Tags = tags
	}

	stats := &models.AccordStatistics{
		Overview:     calculateOverviewStats(accords),
		PyramidStats: calculatePyramidStats(accords),
		TagStats:     calculateTagStats(accords),
		SupplierStats: calculateSupplierStats(accords),
		VolumeStats:  calculateVolumeStats(accords),
		LowInventory: calculateLowInventory(accords, 10.0), // Alert when < 10ml
	}

	return stats, nil
}

func calculateOverviewStats(accords []*models.Accord) models.OverviewStats {
	totalVolume := 0.0
	suppliers := make(map[string]bool)
	tags := make(map[string]bool)

	for _, accord := range accords {
		totalVolume += accord.VolumeMl
		if accord.Supplier != nil && *accord.Supplier != "" {
			suppliers[*accord.Supplier] = true
		}
		for _, tag := range accord.Tags {
			tags[tag] = true
		}
	}

	return models.OverviewStats{
		TotalAccords:   len(accords),
		TotalVolume:    totalVolume,
		TotalSuppliers: len(suppliers),
		TotalTags:      len(tags),
	}
}

func calculatePyramidStats(accords []*models.Accord) models.PyramidStats {
	stats := models.PyramidStats{}

	for _, accord := range accords {
		switch accord.PyramidPosition {
		case "top":
			stats.TopCount++
			stats.TopVolume += accord.VolumeMl
		case "middle":
			stats.MiddleCount++
			stats.MiddleVolume += accord.VolumeMl
		case "base":
			stats.BaseCount++
			stats.BaseVolume += accord.VolumeMl
		}
	}

	return stats
}

func calculateTagStats(accords []*models.Accord) []models.TagStat {
	tagCounts := make(map[string]int)

	for _, accord := range accords {
		for _, tag := range accord.Tags {
			tagCounts[tag]++
		}
	}

	stats := make([]models.TagStat, 0, len(tagCounts))
	for tag, count := range tagCounts {
		stats = append(stats, models.TagStat{
			Tag:   tag,
			Count: count,
		})
	}

	return stats
}

func calculateSupplierStats(accords []*models.Accord) []models.SupplierStat {
	supplierCounts := make(map[string]int)

	for _, accord := range accords {
		if accord.Supplier != nil && *accord.Supplier != "" {
			supplierCounts[*accord.Supplier]++
		}
	}

	stats := make([]models.SupplierStat, 0, len(supplierCounts))
	for supplier, count := range supplierCounts {
		stats = append(stats, models.SupplierStat{
			Supplier: supplier,
			Count:    count,
		})
	}

	return stats
}

func calculateVolumeStats(accords []*models.Accord) models.VolumeStats {
	if len(accords) == 0 {
		return models.VolumeStats{}
	}

	totalVolume := 0.0
	minVolume := accords[0].VolumeMl
	maxVolume := accords[0].VolumeMl

	for _, accord := range accords {
		totalVolume += accord.VolumeMl
		if accord.VolumeMl < minVolume {
			minVolume = accord.VolumeMl
		}
		if accord.VolumeMl > maxVolume {
			maxVolume = accord.VolumeMl
		}
	}

	return models.VolumeStats{
		AverageVolume: totalVolume / float64(len(accords)),
		MinVolume:     minVolume,
		MaxVolume:     maxVolume,
	}
}

func calculateLowInventory(accords []*models.Accord, threshold float64) []models.LowInventoryItem {
	lowItems := make([]models.LowInventoryItem, 0)

	for _, accord := range accords {
		if accord.VolumeMl < threshold {
			lowItems = append(lowItems, models.LowInventoryItem{
				AccordID: accord.ID,
				Name:     accord.Name,
				VolumeML: accord.VolumeMl,
				Supplier: accord.Supplier,
			})
		}
	}

	return lowItems
}
