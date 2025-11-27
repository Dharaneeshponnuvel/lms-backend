package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ContentManagerService struct {
	repo *repositories.ContentManagerRepository
}

func NewContentManagerService(repo *repositories.ContentManagerRepository) *ContentManagerService {
	return &ContentManagerService{repo: repo}
}

func (s *ContentManagerService) CreateWithUser(req models.CreateContentManagerRequest, createdBy uuid.UUID) (uuid.UUID, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return uuid.Nil, errors.New("missing required fields")
	}

	roleID, err := s.repo.GetContentManagerRoleID()
	if err != nil {
		return uuid.Nil, errors.New("content manager role does not exist")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return uuid.Nil, errors.New("failed to hash password")
	}

	user := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		RoleID:       roleID,
		IsActive:     true,
	}

	cm := models.ContentManager{
		CreatedBy: createdBy,
	}

	if err := s.repo.CreateWithUser(&user, &cm); err != nil {
		return uuid.Nil, err
	}

	return cm.UserID, nil
}

func (s *ContentManagerService) AssignBatch(assign *models.ContentManagerBatch) error {
	return s.repo.AssignBatch(assign)
}

func (s *ContentManagerService) GetManagerDetails(userID string) (*models.ContentManager, error) {
	return s.repo.GetByUserID(userID)
}

func (s *ContentManagerService) GetManagerBatches(managerID string) ([]models.ContentManagerBatch, error) {
	return s.repo.GetAssignedBatches(managerID)
}
