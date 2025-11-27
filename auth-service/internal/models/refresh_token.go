package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"sessionId"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
	IsUsed    bool      `gorm:"default:false" json:"isUsed"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

	User    User    `gorm:"foreignKey:UserID"`
	Session Session `gorm:"foreignKey:SessionID"`
}
