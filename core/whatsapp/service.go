package whatsapp

import (
	"context"
	"time"

	cmn "api_go/core/common"
)

type whatsAppServiceImpl struct {
	repo     WhatsAppRepository
	provider WhatsAppProvider // Ahora está definido
}

// NewWhatsAppService crea una nueva instancia del servicio
func NewWhatsAppService(repo WhatsAppRepository, provider WhatsAppProvider) WhatsAppService {
	return &whatsAppServiceImpl{
		repo:     repo,
		provider: provider,
	}
}

func (s *whatsAppServiceImpl) RegisterSession(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Validar que la sesión no exista
	existing, _ := s.repo.FindSession(ctx, req.SessionID)
	if existing != nil {
		return nil, cmn.NewValidationError("session_id", "La sesión ya existe")
	}

	// Crear nueva sesión
	session := &Session{
		ID:         req.SessionID,
		Provider:   req.Provider,
		Status:     StatusPending,
		Metadata:   req.Metadata,
		WebhookURL: req.WebhookURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Registrar en el proveedor externo
	providerResp, err := s.provider.Register(ctx, session)
	if err != nil {
		return nil, cmn.NewServiceError("Error registrando sesión en proveedor: " + err.Error())
	}

	// Actualizar sesión con datos del proveedor
	session.QRCode = providerResp.QRCode
	session.Status = providerResp.Status

	// Guardar en repositorio
	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		QRCode:    session.QRCode,
		Status:    session.Status,
		ExpiresAt: providerResp.ExpiresAt,
	}, nil
}

// CORREGIDO: Ahora retorna *Session en lugar de *SessionStatus
func (s *whatsAppServiceImpl) GetSessionStatus(ctx context.Context, sessionID string) (*Session, error) {
	session, err := s.repo.FindSession(ctx, sessionID)
	if err != nil {
		return nil, cmn.NewNotFoundError("Sesión no encontrada")
	}
	if session == nil {
		return nil, cmn.NewNotFoundError("Sesión no encontrada")
	}

	// Obtener estado actual del proveedor
	status, err := s.provider.GetStatus(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Actualizar estado en repositorio si cambió
	if session.Status != status {
		session.Status = status
		session.UpdatedAt = time.Now()
		if status == StatusConnected && session.ConnectedAt == nil {
			now := time.Now()
			session.ConnectedAt = &now
		}
		s.repo.SaveSession(ctx, session)
	}

	return session, nil
}

func (s *whatsAppServiceImpl) SendMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	session, err := s.repo.FindSession(ctx, req.SessionID)
	if err != nil {
		return nil, cmn.NewNotFoundError("Sesión no encontrada")
	}
	if session == nil {
		return nil, cmn.NewNotFoundError("Sesión no encontrada")
	}

	if session.Status != StatusConnected {
		return nil, cmn.NewValidationError("session", "La sesión no está conectada")
	}

	// Enviar mensaje a través del proveedor
	result, err := s.provider.SendMessage(ctx, req)
	if err != nil {
		return nil, cmn.NewServiceError("Error enviando mensaje: " + err.Error())
	}

	return &MessageResponse{
		Data: result,
	}, nil
}

func (s *whatsAppServiceImpl) DeleteSession(ctx context.Context, sessionID string) error {
	// Eliminar del proveedor
	if err := s.provider.Unregister(ctx, sessionID); err != nil {
		return err
	}

	// Eliminar del repositorio
	return s.repo.DeleteSession(ctx, sessionID)
}

func (s *whatsAppServiceImpl) ListSessions(ctx context.Context) ([]*Session, error) {
	return s.repo.ListSessions(ctx)
}
