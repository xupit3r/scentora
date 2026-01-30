package services

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/yourusername/scentora-backend/internal/models"
	"github.com/yourusername/scentora-backend/internal/repository"
)

type PerfumeService struct {
	repo *repository.PerfumeRepository
}

func NewPerfumeService(repo *repository.PerfumeRepository) *PerfumeService {
	return &PerfumeService{repo: repo}
}

func (s *PerfumeService) Create(userID string, req *models.CreatePerfumeRequest) (*models.PerfumeResponse, error) {
	perfume := &models.Perfume{
		UserID:        userID,
		Name:          req.Name,
		Designer:      req.Designer,
		Year:          req.Year,
		Concentration: req.Concentration,
		TopNotes:      req.Pyramid.Top,
		MiddleNotes:   req.Pyramid.Middle,
		BaseNotes:     req.Pyramid.Base,
		Description:   req.Description,
		ImageURL:      req.ImageURL,
	}

	if err := s.repo.Create(perfume); err != nil {
		return nil, err
	}

	return perfume.ToResponse(), nil
}

func (s *PerfumeService) Get(id, userID string) (*models.PerfumeResponse, error) {
	perfume, err := s.repo.FindByID(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("perfume not found")
		}
		return nil, err
	}
	return perfume.ToResponse(), nil
}

func (s *PerfumeService) List(userID string, filters map[string]string) ([]*models.PerfumeResponse, error) {
	perfumes, err := s.repo.List(userID, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.PerfumeResponse, len(perfumes))
	for i, p := range perfumes {
		responses[i] = p.ToResponse()
	}
	return responses, nil
}

func (s *PerfumeService) Update(id, userID string, req *models.UpdatePerfumeRequest) (*models.PerfumeResponse, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Designer != nil {
		updates["designer"] = *req.Designer
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.Concentration != nil {
		updates["concentration"] = *req.Concentration
	}
	if req.Pyramid != nil {
		updates["top_notes"] = pq.Array(req.Pyramid.Top)
		updates["middle_notes"] = pq.Array(req.Pyramid.Middle)
		updates["base_notes"] = pq.Array(req.Pyramid.Base)
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}

	if len(updates) == 0 {
		return s.Get(id, userID)
	}

	perfume, err := s.repo.Update(id, userID, updates)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("perfume not found")
		}
		return nil, err
	}

	return perfume.ToResponse(), nil
}

func (s *PerfumeService) Delete(id, userID string) error {
	err := s.repo.Delete(id, userID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("perfume not found")
	}
	return err
}
