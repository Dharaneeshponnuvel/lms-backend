package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Batch struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BatchYearID uuid.UUID      `gorm:"type:uuid;not null;index" json:"batch_year_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations (ignored in requests)
	BatchYear BatchYear `gorm:"foreignKey:BatchYearID" json:"batch_year,omitempty"`
	Creator   User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Students  []Student `gorm:"foreignKey:BatchID" json:"students,omitempty"`
}
