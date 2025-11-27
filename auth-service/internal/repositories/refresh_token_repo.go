package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db}
}

func (r *RefreshTokenRepository) Create(rt *models.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *RefreshTokenRepository) MarkUsed(id string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("id = ?", id).Update("is_used", true).Error
}
