package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentManager struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"`
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User    User                  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Creator User                  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Batches []ContentManagerBatch `gorm:"foreignKey:ContentManagerID" json:"batches,omitempty"`
}

type CreateContentManagerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}
