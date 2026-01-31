package services

import (
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type TagService struct {
	tagRepo *repository.PredefinedTagRepository
}

func NewTagService(tagRepo *repository.PredefinedTagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// GetAllTags retrieves all predefined tags
func (s *TagService) GetAllTags() ([]*models.PredefinedTag, error) {
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}
	return tags, nil
}

// GetTagsByCategory retrieves tags for a specific category
func (s *TagService) GetTagsByCategory(category string) ([]*models.PredefinedTag, error) {
	tags, err := s.tagRepo.GetByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by category: %w", err)
	}
	return tags, nil
}

// SearchTags searches for tags by partial match
func (s *TagService) SearchTags(search string) ([]*models.PredefinedTag, error) {
	if search == "" {
		return []*models.PredefinedTag{}, nil
	}

	tags, err := s.tagRepo.Search(search)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}
	return tags, nil
}

// GetAllCategories retrieves all unique categories
func (s *TagService) GetAllCategories() ([]string, error) {
	categories, err := s.tagRepo.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// GetTagsGroupedByCategory retrieves all tags grouped by category
func (s *TagService) GetTagsGroupedByCategory() (map[string][]string, error) {
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}

	grouped := make(map[string][]string)
	for _, tag := range tags {
		grouped[tag.Category] = append(grouped[tag.Category], tag.Tag)
	}

	return grouped, nil
}
