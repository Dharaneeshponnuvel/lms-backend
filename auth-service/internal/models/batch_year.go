package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BatchYear struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InstitutionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"institution_id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"` // 2024–2025
	StartYear     time.Time      `gorm:"not null;index" json:"start_year"`
	EndYear       time.Time      `gorm:"not null" json:"end_year"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations (ignored in POST requests)
	Institution Institution `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`
	Creator     User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Batches     []Batch     `gorm:"foreignKey:BatchYearID" json:"batches,omitempty"`
}
