package repositories

import (
	"auth-service/internal/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ContentManagerRepository struct {
	db *gorm.DB
}

func NewContentManagerRepository(db *gorm.DB) *ContentManagerRepository {
	return &ContentManagerRepository{db: db}
}

// ➕ Create Content Manager User
func (r *ContentManagerRepository) Create(cm *models.ContentManager) error {
	return r.db.Create(cm).Error
}

// 🔗 Assign Content Manager to Batch
func (r *ContentManagerRepository) AssignBatch(assign *models.ContentManagerBatch) error {
	return r.db.Create(assign).Error
}

// 🔍 Get all batches managed by Content Manager
func (r *ContentManagerRepository) GetAssignedBatches(managerID string) ([]models.ContentManagerBatch, error) {
	var assignments []models.ContentManagerBatch
	err := r.db.Where("content_manager_id = ?", managerID).Find(&assignments).Error
	return assignments, err
}

// 🔍 Get Content Manager Profile
func (r *ContentManagerRepository) GetByUserID(userID string) (*models.ContentManager, error) {
	var cm models.ContentManager
	err := r.db.Where("user_id = ?", userID).First(&cm).Error
	return &cm, err
}
func (r *ContentManagerRepository) CreateWithUser(user *models.User, cm *models.ContentManager) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		cm.UserID = user.ID
		return tx.Create(cm).Error
	})
}
func (r *ContentManagerRepository) GetContentManagerRoleID() (uuid.UUID, error) {
	var role models.Role
	err := r.db.Where("name = ?", "CONTENT_MANAGER").First(&role).Error
	return role.ID, err
}
