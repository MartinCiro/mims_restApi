package external

import (
	"context"
	core_wsp "scrapper_go_email/core/whatsapp"
	"time"
)

type WhatsMeowAdapter struct {
	// Aquí irían las dependencias de WhatsMeow
}

func NewWhatsMeowAdapter() *WhatsMeowAdapter {
	return &WhatsMeowAdapter{}
}

func (a *WhatsMeowAdapter) Register(ctx context.Context, session *core_wsp.Session) (*core_wsp.RegisterResponse, error) {
	// Implementación mock por ahora
	expiresAt := time.Now().Add(5 * time.Minute)
	return &core_wsp.RegisterResponse{
		QRCode:    "mock_qr_code_base64",
		Status:    core_wsp.StatusPending,
		ExpiresAt: &expiresAt,
	}, nil
}

func (a *WhatsMeowAdapter) GetStatus(ctx context.Context, sessionID string) (core_wsp.SessionStatus, error) {
	return core_wsp.StatusPending, nil
}

func (a *WhatsMeowAdapter) SendMessage(ctx context.Context, req *core_wsp.MessageRequest) (*core_wsp.MessageResult, error) {
	return &core_wsp.MessageResult{
		MessageID: "mock_message_id",
		To:        req.To,
		SentAt:    time.Now(),
		Status:    "sent",
	}, nil
}

func (a *WhatsMeowAdapter) Unregister(ctx context.Context, sessionID string) error {
	return nil
}
