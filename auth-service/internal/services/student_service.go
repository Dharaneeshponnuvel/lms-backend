package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type StudentService struct {
	repo *repositories.StudentRepository
}

func NewStudentService(repo *repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// Create Student
func (s *StudentService) CreateStudent(student *models.Student) error {
	return s.repo.Create(student)
}

// Get students by batch
func (s *StudentService) GetByBatch(batchID string) ([]models.Student, error) {
	return s.repo.FindByBatch(batchID)
}
