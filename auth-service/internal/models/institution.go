package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Institution struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"` // Institution Admin
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	AdminPosition string         `gorm:"type:varchar(100)" json:"admin_position"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Settings      datatypes.JSON `gorm:"type:jsonb" json:"settings,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User        `gorm:"foreignKey:UserID" json:"user"`
	BatchYears []BatchYear `gorm:"foreignKey:InstitutionID" json:"batch_years"`
}
