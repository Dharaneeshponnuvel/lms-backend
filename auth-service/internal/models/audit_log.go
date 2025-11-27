package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID     `gorm:"type:uuid"`
	Action    string         `gorm:"type:varchar(100)"`
	Route     string         `gorm:"type:varchar(255)"`
	IP        string         `gorm:"type:varchar(50)"`
	UserAgent string         `gorm:"type:text"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
