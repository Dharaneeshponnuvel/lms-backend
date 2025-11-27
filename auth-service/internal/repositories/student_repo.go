package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type StudentRepository struct{ db *gorm.DB }

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db}
}

func (r *StudentRepository) Create(s *models.Student) error {
	return r.db.Create(s).Error
}

func (r *StudentRepository) FindByBatch(batchID string) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Where("batch_id = ?", batchID).Preload("User").Find(&students).Error
	return students, err
}
