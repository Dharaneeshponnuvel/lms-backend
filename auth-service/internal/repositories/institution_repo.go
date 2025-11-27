package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type InstitutionRepository struct{ db *gorm.DB }

func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{db}
}

func (r *InstitutionRepository) Create(inst *models.Institution) error {
	return r.db.Create(inst).Error
}

func (r *InstitutionRepository) FindByID(id string) (*models.Institution, error) {
	var inst models.Institution
	err := r.db.First(&inst, "id = ?", id).Error
	return &inst, err
}

func (r *InstitutionRepository) FindByUserID(userID string) (*models.Institution, error) {
	var inst models.Institution
	err := r.db.Where("user_id = ?", userID).First(&inst).Error
	return &inst, err
}

func (r *InstitutionRepository) GetAll() ([]models.Institution, error) {
	var insts []models.Institution
	err := r.db.Find(&insts).Error
	return insts, err
}

func (r *InstitutionRepository) Update(inst *models.Institution) error {
	return r.db.Save(inst).Error
}

func (r *InstitutionRepository) Delete(id string) error {
	return r.db.Delete(&models.Institution{}, "id = ?", id).Error
}
