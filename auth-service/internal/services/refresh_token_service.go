package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type RefreshTokenService struct {
	repo *repositories.RefreshTokenRepository
}

func NewRefreshTokenService(repo *repositories.RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{repo: repo}
}

func (s *RefreshTokenService) SaveToken(rt *models.RefreshToken) error {
	return s.repo.Create(rt)
}

func (s *RefreshTokenService) MarkUsed(id string) error {
	return s.repo.MarkUsed(id)
}
