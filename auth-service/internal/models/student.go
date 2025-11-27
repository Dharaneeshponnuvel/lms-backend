package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Student struct {
	UserID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"user_id"`
	BatchID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"batch_id"`
	RollNumber      string         `gorm:"type:varchar(50);uniqueIndex:idx_batch_roll" json:"roll_number"`
	Branch          string         `gorm:"type:varchar(100)" json:"branch"`
	DateOfAdmission time.Time      `json:"date_of_admission"` // remove or handle
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	Demographics    datatypes.JSON `gorm:"type:jsonb" json:"demographics"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
