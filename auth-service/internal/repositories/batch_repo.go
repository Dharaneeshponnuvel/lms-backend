package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type BatchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db}
}

func (r *BatchRepository) Create(batch *models.Batch) error {
	return r.db.Create(batch).Error
}

func (r *BatchRepository) GetByBatchYear(batchYearID string) ([]models.Batch, error) {
	var batches []models.Batch
	err := r.db.Where("batch_year_id = ?", batchYearID).Find(&batches).Error
	return batches, err
}

func (r *BatchRepository) GetByID(id string) (*models.Batch, error) {
	var batch models.Batch
	err := r.db.First(&batch, "id = ?", id).Error
	return &batch, err
}
