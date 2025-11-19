package redis

import (
	"context"
	"fmt"
	"time"
)

type SessionManager struct {
	cache            *Cache
	sessionTTL       time.Duration
	refreshThreshold time.Duration
}

func NewSessionManager(cache *Cache, sessionTTL, refreshThreshold time.Duration) *SessionManager {
	return &SessionManager{
		cache:            cache,
		sessionTTL:       sessionTTL,
		refreshThreshold: refreshThreshold,
	}
}

// CreateSession crea una nueva sesión
func (sm *SessionManager) CreateSession(ctx context.Context, sessionKey string, sessionData map[string]interface{}) error {
	return sm.cache.Set(ctx, sessionKey, sessionData, sm.sessionTTL)
}

// GetSession obtiene los datos de una sesión
func (sm *SessionManager) GetSession(ctx context.Context, sessionKey string) (map[string]interface{}, error) {
	var sessionData map[string]interface{}
	err := sm.cache.GetJSON(ctx, sessionKey, &sessionData)
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}

// CheckAndRefreshSession verifica si la sesión necesita refresh y la refresca
func (sm *SessionManager) CheckAndRefreshSession(ctx context.Context, sessionKey string) (bool, error) {
	// Verificar si está por expirar
	needsRefresh, _, err := sm.cache.IsAboutToExpire(ctx, sessionKey, sm.refreshThreshold)
	if err != nil {
		return false, err
	}

	if needsRefresh {
		// Refrescar la sesión
		err = sm.cache.RefreshTTL(ctx, sessionKey, sm.sessionTTL)
		if err != nil {
			return false, fmt.Errorf("error refrescando sesión: %v", err)
		}
		return true, nil
	}

	return false, nil
}

// GetSessionInfo obtiene información completa de la sesión
func (sm *SessionManager) GetSessionInfo(ctx context.Context, sessionKey string) (map[string]interface{}, error) {
	info, err := sm.cache.GetTTLInfo(ctx, sessionKey)
	if err != nil {
		return nil, err
	}

	// Agregar datos de la sesión si existe
	sessionData, err := sm.GetSession(ctx, sessionKey)
	if err == nil && sessionData != nil {
		info["session_data"] = sessionData
	}

	// Agregar información de refresh
	needsRefresh, _, _ := sm.cache.IsAboutToExpire(ctx, sessionKey, sm.refreshThreshold)
	info["needs_refresh"] = needsRefresh
	info["refresh_threshold_seconds"] = int64(sm.refreshThreshold.Seconds())

	return info, nil
}

// DeleteSession elimina una sesión
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionKey string) error {
	return sm.cache.Delete(ctx, sessionKey)
}

// SessionExists verifica si una sesión existe
func (sm *SessionManager) SessionExists(ctx context.Context, sessionKey string) (bool, error) {
	return sm.cache.Exists(ctx, sessionKey)
}
