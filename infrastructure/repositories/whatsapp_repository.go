package repositories

import (
	core_wsp "api_go/core/whatsapp"
	"context"
	"sync"
)

type WhatsAppMemoryRepository struct {
	sessions map[string]*core_wsp.Session
	mu       sync.RWMutex
}

func NewWhatsAppMemoryRepository() *WhatsAppMemoryRepository {
	return &WhatsAppMemoryRepository{
		sessions: make(map[string]*core_wsp.Session),
	}
}

func (r *WhatsAppMemoryRepository) SaveSession(ctx context.Context, session *core_wsp.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *WhatsAppMemoryRepository) FindSession(ctx context.Context, sessionID string) (*core_wsp.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, nil
	}
	return session, nil
}

func (r *WhatsAppMemoryRepository) DeleteSession(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, sessionID)
	return nil
}

func (r *WhatsAppMemoryRepository) ListSessions(ctx context.Context) ([]*core_wsp.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]*core_wsp.Session, 0, len(r.sessions))
	for _, session := range r.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *WhatsAppMemoryRepository) UpdateSessionStatus(ctx context.Context, sessionID string, status core_wsp.SessionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if session, exists := r.sessions[sessionID]; exists {
		session.Status = status
	}
	return nil
}
