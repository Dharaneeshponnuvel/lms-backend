package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type InstitutionService struct {
	repo *repositories.InstitutionRepository
}

func NewInstitutionService(repo *repositories.InstitutionRepository) *InstitutionService {
	return &InstitutionService{repo}
}

func (s *InstitutionService) Create(inst *models.Institution) error {
	return s.repo.Create(inst)
}

func (s *InstitutionService) GetByID(id string) (*models.Institution, error) {
	return s.repo.FindByID(id)
}

func (s *InstitutionService) GetByUserID(userID string) (*models.Institution, error) {
	return s.repo.FindByUserID(userID)
}
