package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type BatchYearService struct {
	repo *repositories.BatchYearRepository
}

// 💡 Always inject repository, not DB
func NewBatchYearService(repo *repositories.BatchYearRepository) *BatchYearService {
	return &BatchYearService{repo: repo}
}

func (s *BatchYearService) Create(batchYear *models.BatchYear) error {
	return s.repo.Create(batchYear)
}

func (s *BatchYearService) GetInstitutionBatchYears(instID string) ([]models.BatchYear, error) {
	return s.repo.GetByInstitution(instID)
}

func (s *BatchYearService) GetByID(id string) (*models.BatchYear, error) {
	return s.repo.GetByID(id)
}
