package services

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type JournalService struct {
	journalRepo *repository.JournalRepository
	perfumeRepo *repository.PerfumeRepository
}

func NewJournalService(journalRepo *repository.JournalRepository, perfumeRepo *repository.PerfumeRepository) *JournalService {
	return &JournalService{
		journalRepo: journalRepo,
		perfumeRepo: perfumeRepo,
	}
}

func (s *JournalService) Create(userID string, req *models.CreateJournalRequest) (*models.JournalEntry, error) {
	// Verify perfume exists and belongs to user
	_, err := s.perfumeRepo.FindByID(req.PerfumeID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("perfume not found")
		}
		return nil, err
	}

	entry := &models.JournalEntry{
		UserID:    userID,
		PerfumeID: req.PerfumeID,
		Date:      req.Date,
		Content:   req.Content,
		Rating:    req.Rating,
		Occasion:  req.Occasion,
		Weather:   req.Weather,
	}

	if err := s.journalRepo.Create(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *JournalService) ListByPerfume(perfumeID, userID string) ([]*models.JournalEntry, error) {
	// Verify perfume exists and belongs to user
	_, err := s.perfumeRepo.FindByID(perfumeID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("perfume not found")
		}
		return nil, err
	}

	return s.journalRepo.ListByPerfume(perfumeID, userID)
}

func (s *JournalService) Update(id, userID string, req *models.UpdateJournalRequest) (*models.JournalEntry, error) {
	updates := make(map[string]interface{})

	if req.Date != nil {
		updates["date"] = *req.Date
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Occasion != nil {
		updates["occasion"] = *req.Occasion
	}
	if req.Weather != nil {
		updates["weather"] = *req.Weather
	}

	if len(updates) == 0 {
		return s.journalRepo.FindByID(id, userID)
	}

	entry, err := s.journalRepo.Update(id, userID, updates)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, err
	}

	return entry, nil
}

func (s *JournalService) Delete(id, userID string) error {
	err := s.journalRepo.Delete(id, userID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("journal entry not found")
	}
	return err
}
