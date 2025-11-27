package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Session struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Token        string         `gorm:"type:text;not null;uniqueIndex" json:"token"`
	RefreshToken string         `gorm:"type:text;uniqueIndex" json:"refreshToken"`
	DeviceInfo   datatypes.JSON `gorm:"type:jsonb" json:"deviceInfo"`
	ExpiresAt    time.Time      `gorm:"not null;index" json:"expiresAt"`
	IsActive     bool           `gorm:"default:true;index" json:"isActive"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationship
	User User `gorm:"foreignKey:UserID"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
