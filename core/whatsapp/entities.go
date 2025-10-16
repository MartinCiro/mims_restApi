package whatsapp

import (
	"time"
)

type SessionStatus string

const (
	StatusPending      SessionStatus = "pending"
	StatusConnected    SessionStatus = "connected"
	StatusDisconnected SessionStatus = "disconnected"
	StatusFailed       SessionStatus = "failed"
)

type RegisterRequest struct {
	SessionID  string            `json:"session_id" validate:"required"`
	Provider   string            `json:"provider" validate:"required"`
	WebhookURL string            `json:"webhook_url,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// RegisterResponse ahora es una estructura separada
type RegisterResponse struct {
	QRCode    string        `json:"qr_code,omitempty"`
	Status    SessionStatus `json:"status"`
	ExpiresAt *time.Time    `json:"expires_at,omitempty"`
}

type MessageRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	To        string `json:"to" validate:"required"`
	Message   string `json:"message" validate:"required"`
	Type      string `json:"type" validate:"required"`
}

type MessageResponse struct {
	Data *MessageResult `json:"data,omitempty"`
}

type Session struct {
	ID          string            `json:"id"`
	Provider    string            `json:"provider"`
	Status      SessionStatus     `json:"status"`
	QRCode      string            `json:"qr_code,omitempty"`
	UserInfo    *UserInfo         `json:"user_info,omitempty"`
	WebhookURL  string            `json:"webhook_url,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ConnectedAt *time.Time        `json:"connected_at,omitempty"`
}

type UserInfo struct {
	Phone      string `json:"phone"`
	Name       string `json:"name,omitempty"`
	ProfilePic string `json:"profile_pic,omitempty"`
	Status     string `json:"status,omitempty"`
}

type MessageResult struct {
	MessageID string    `json:"message_id"`
	To        string    `json:"to"`
	SentAt    time.Time `json:"sent_at"`
	Status    string    `json:"status"`
}
