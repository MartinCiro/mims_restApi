package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID          string         `gorm:"primaryKey;size:100" json:"id"`
	Provider    string         `gorm:"size:50;not null" json:"provider"`
	Status      string         `gorm:"size:20;not null" json:"status"`
	QRCode      string         `gorm:"type:text" json:"qr_code,omitempty"`
	WebhookURL  string         `gorm:"size:500" json:"webhook_url,omitempty"`
	Metadata    string         `gorm:"type:text" json:"metadata"`  // JSON como string
	UserInfo    string         `gorm:"type:text" json:"user_info"` // JSON como string
	ConnectedAt *time.Time     `json:"connected_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla
func (Session) TableName() string {
	return "whatsapp_sessions"
}
