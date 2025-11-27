package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type SessionService struct {
	repo *repositories.SessionRepository
}

func NewSessionService(repo *repositories.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(session *models.Session) error {
	return s.repo.Create(session)
}

// 🧠 ADD THIS METHOD — VERY IMPORTANT ❗
func (s *SessionService) InvalidateOldSessions(userID string) {
	s.repo.InvalidateAll(userID)
}

func (s *SessionService) GetActiveSessions(userID string) ([]models.Session, error) {
	return s.repo.GetActiveSessions(userID)
}

func (s *SessionService) DeactivateSession(userID, token string) error {
	return s.repo.DeactivateSession(userID, token)
}
