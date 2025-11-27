package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type BatchService struct {
	repo *repositories.BatchRepository
}

func NewBatchService(repo *repositories.BatchRepository) *BatchService {
	return &BatchService{repo: repo}
}

func (s *BatchService) Create(batch *models.Batch) error {
	return s.repo.Create(batch)
}

func (s *BatchService) GetByBatchYear(batchYearID string) ([]models.Batch, error) {
	return s.repo.GetByBatchYear(batchYearID)
}
