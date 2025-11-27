package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

// ➕ Create role
func (s *RoleService) CreateRole(role *models.Role) error {
	return s.repo.Create(role)
}

// 🔍 Get all roles
func (s *RoleService) GetRoles() ([]models.Role, error) {
	return s.repo.GetAll()
}

// 🔍 Find by name
func (s *RoleService) FindByName(name string) (*models.Role, error) {
	return s.repo.FindByName(name)
}

// ❌ Delete role
func (s *RoleService) DeleteRole(id string) error {
	return s.repo.Delete(id)
}
