package models

import (
	"time"

	"github.com/google/uuid"
)

type ContentManagerBatch struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ContentManagerID uuid.UUID `gorm:"type:uuid;not null;index" json:"content_manager_id"`
	BatchID          uuid.UUID `gorm:"type:uuid;not null;index" json:"batch_id"`
	AssignedBy       uuid.UUID `gorm:"type:uuid;not null" json:"assigned_by"`
	AssignedAt       time.Time `gorm:"autoCreateTime" json:"assigned_at"`

	// Relations (omit in API response using `omitempty`)
	ContentManager ContentManager `gorm:"foreignKey:ContentManagerID" json:"content_manager,omitempty"`
	Batch          Batch          `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	Assigner       User           `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
}
