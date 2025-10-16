package whatsapp

import (
	"context"
)

// WhatsAppService port (interfaz)
type WhatsAppService interface {
	RegisterSession(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	GetSessionStatus(ctx context.Context, sessionID string) (*Session, error) // Cambiado a *Session
	SendMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error)
	DeleteSession(ctx context.Context, sessionID string) error
	ListSessions(ctx context.Context) ([]*Session, error)
}

// WhatsAppRepository port (para persistencia)
type WhatsAppRepository interface {
	SaveSession(ctx context.Context, session *Session) error
	FindSession(ctx context.Context, sessionID string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	ListSessions(ctx context.Context) ([]*Session, error)
	UpdateSessionStatus(ctx context.Context, sessionID string, status SessionStatus) error
}

// WhatsAppProvider port para proveedores externos (NUEVO)
type WhatsAppProvider interface {
	Register(ctx context.Context, session *Session) (*RegisterResponse, error)
	GetStatus(ctx context.Context, sessionID string) (SessionStatus, error)
	SendMessage(ctx context.Context, req *MessageRequest) (*MessageResult, error)
	Unregister(ctx context.Context, sessionID string) error
}
