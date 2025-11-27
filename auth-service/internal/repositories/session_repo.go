package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type SessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) InvalidateAll(userID string) {
	r.db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Update("is_active", false)
}

func (r *SessionRepository) GetActiveSessions(userID string) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Where("user_id = ? AND is_active = true", userID).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) DeactivateSession(userID, token string) error {
	return r.db.Model(&models.Session{}).
		Where("user_id = ? AND token = ?", userID, token).
		Update("is_active", false).Error
}
