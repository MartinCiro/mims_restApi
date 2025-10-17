package external

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"api_go/core/whatsapp"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsMeowAdapter struct {
	client      *whatsmeow.Client
	container   *sqlstore.Container
	qrChan      <-chan whatsmeow.QRChannelItem
	currentQR   string
	isConnected bool
	qrMutex     sync.RWMutex
	sessionID   string
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWhatsMeowAdapter() *WhatsMeowAdapter {
	ctx, cancel := context.WithCancel(context.Background())
	return &WhatsMeowAdapter{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *WhatsMeowAdapter) Register(ctx context.Context, session *whatsapp.Session) (*whatsapp.RegisterResponse, error) {
	a.sessionID = session.ID

	log.Printf("🚀 Iniciando nueva sesión de WhatsApp para %s...", session.ID)

	// Configuración del logger
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Configuración de la base de datos (SQLite)
	container, err := sqlstore.New(a.ctx, "sqlite",
		fmt.Sprintf("file:%s_whatsapp.db?_pragma=foreign_keys=1&_pragma=journal_mode=WAL&_pragma=busy_timeout=10000&_pragma=synchronous=NORMAL&_pragma=cache_size=10000", session.ID),
		dbLog)
	if err != nil {
		return nil, fmt.Errorf("error creando almacenamiento: %v", err)
	}

	a.container = container

	// Si no hay dispositivos guardados, registrar uno nuevo
	deviceStore, err := container.GetFirstDevice(a.ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo dispositivo: %v", err)
	}

	// Configuración del cliente
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	a.client = client

	// Registrar handler de eventos simples
	a.client.AddEventHandler(a.eventHandler)

	// Conectar cliente
	if a.client.Store.ID == nil {
		// No hay sesión guardada, necesitamos iniciar sesión
		log.Printf("📱 No hay sesión guardada, generando QR code...")

		// Obtener canal QR
		qrChan, err := a.client.GetQRChannel(a.ctx)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo canal QR: %v", err)
		}
		a.qrChan = qrChan

		// Manejar el QR en una goroutine (igual que tu código original)
		go a.handleQRChannel()

		// Conectar el cliente
		err = a.client.Connect()
		if err != nil {
			return nil, fmt.Errorf("error conectando: %v", err)
		}

		// Esperar un momento para que el QR se genere
		time.Sleep(2 * time.Second)

		// Retornar el QR code generado
		if a.currentQR != "" {
			qrDataURL, err := a.generateQRImage(a.currentQR)
			if err != nil {
				return nil, fmt.Errorf("error generando QR image: %v", err)
			}

			expiresAt := time.Now().Add(5 * time.Minute)
			return &whatsapp.RegisterResponse{
				QRCode:    qrDataURL,
				Status:    whatsapp.StatusPending,
				ExpiresAt: &expiresAt,
			}, nil
		}

		// Si no hay QR aún, retornar estado pending
		return &whatsapp.RegisterResponse{
			Status: whatsapp.StatusPending,
		}, nil

	} else {
		// Sesión ya guardada
		log.Printf("🔗 Sesión existente encontrada para %s", session.ID)

		// Conectar directamente
		err = a.client.Connect()
		if err != nil {
			return nil, fmt.Errorf("error conectando: %v", err)
		}

		a.setConnected(true)

		expiresAt := time.Now().Add(24 * time.Hour)
		return &whatsapp.RegisterResponse{
			Status:    whatsapp.StatusConnected,
			ExpiresAt: &expiresAt,
		}, nil
	}
}

// handleQRChannel maneja el canal QR (igual que tu código original)
func (a *WhatsMeowAdapter) handleQRChannel() {
	for evt := range a.qrChan {
		switch evt.Event {
		case "code":
			log.Printf("📱 Nuevo código QR generado para sesión %s", a.sessionID)
			a.setQRCode(evt.Code)

		case "success":
			log.Printf("✅ ¡Sesión %s iniciada correctamente!", a.sessionID)
			a.setConnected(true)
			a.clearQRCode()

		case "timeout":
			log.Printf("❌ Tiempo agotado para sesión %s", a.sessionID)
			a.clearQRCode()
			// Reconectar automáticamente después de timeout
			go a.reconnectAfterTimeout()
		}
	}
}

// reconnectAfterTimeout reconecta después de un timeout
func (a *WhatsMeowAdapter) reconnectAfterTimeout() {
	time.Sleep(2 * time.Second)
	log.Printf("🔄 Reconectando después de timeout...")

	// Desconectar si está conectado
	if a.client != nil && a.client.IsConnected() {
		a.client.Disconnect()
	}

	// Reconectar
	if err := a.client.Connect(); err != nil {
		log.Printf("❌ Error reconectando: %v", err)
	}
}

// eventHandler maneja eventos de conexión
func (a *WhatsMeowAdapter) eventHandler(evt interface{}) {
	switch evt.(type) {
	case *events.Connected:
		log.Printf("🔗 WhatsApp conectado para sesión %s", a.sessionID)
		a.setConnected(true)

	case *events.Disconnected:
		log.Printf("🔌 WhatsApp desconectado para sesión %s", a.sessionID)
		a.setConnected(false)

	case *events.LoggedOut:
		log.Printf("🚪 Sesión %s cerrada", a.sessionID)
		a.setConnected(false)
	}
}

func (a *WhatsMeowAdapter) GetStatus(ctx context.Context, sessionID string) (whatsapp.SessionStatus, error) {
	if a.client == nil {
		return whatsapp.StatusDisconnected, nil
	}

	if a.isConnected || (a.client.IsConnected() && a.client.IsLoggedIn()) {
		return whatsapp.StatusConnected, nil
	}

	if a.currentQR != "" {
		return whatsapp.StatusPending, nil
	}

	return whatsapp.StatusDisconnected, nil
}

func (a *WhatsMeowAdapter) SendMessage(ctx context.Context, req *whatsapp.MessageRequest) (*whatsapp.MessageResult, error) {
	if a.client == nil || !a.client.IsConnected() || !a.isConnected {
		return nil, fmt.Errorf("cliente no conectado")
	}

	// Validar y formatear el número de teléfono
	recipient, err := a.validatePhoneNumber(req.To)
	if err != nil {
		return nil, fmt.Errorf("número de teléfono inválido: %v", err)
	}

	// Enviar mensaje usando la misma estructura que tu código original
	msg := &waE2E.Message{
		Conversation: proto.String(req.Message),
	}

	resp, err := a.client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return nil, fmt.Errorf("error enviando mensaje: %v", err)
	}

	log.Printf("✅ Mensaje enviado a %s", req.To)

	return &whatsapp.MessageResult{
		MessageID: resp.ID,
		To:        req.To,
		SentAt:    time.Now(),
		Status:    "sent",
	}, nil
}

func (a *WhatsMeowAdapter) Unregister(ctx context.Context, sessionID string) error {
	if a.cancel != nil {
		a.cancel()
	}

	if a.client != nil {
		if a.client.IsConnected() {
			a.client.Disconnect()
		}
		a.client = nil
	}
	a.clearQRCode()
	a.setConnected(false)

	log.Printf("🔴 Sesión %s desconectada", sessionID)
	return nil
}

// validatePhoneNumber - CORREGIDO para retornar types.JID
func (a *WhatsMeowAdapter) validatePhoneNumber(phone string) (types.JID, error) {
	// Limpiar el número (eliminar espacios, guiones, etc.)
	cleanPhone := strings.TrimSpace(phone)
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")

	var formattedPhone string

	// Si el número ya tiene código de país (+57, +52, +1, etc.)
	if strings.HasPrefix(cleanPhone, "+") {
		// Remover el + y usar el número tal cual
		cleanPhone = strings.TrimPrefix(cleanPhone, "+")
		formattedPhone = fmt.Sprintf("%s@s.whatsapp.net", cleanPhone)
	} else {
		// Si no tiene código de país, asumir que es un número local
		// y agregar código por defecto (puedes cambiar este según tu necesidad)
		// Por ejemplo: México = 52, Colombia = 57, Argentina = 54, etc.
		defaultCountryCode := "57" // Cambia este según tu país por defecto
		formattedPhone = fmt.Sprintf("%s%s@s.whatsapp.net", defaultCountryCode, cleanPhone)
	}

	// Validar que el número tenga al menos 10 dígitos (código país + número)
	if len(cleanPhone) < 10 {
		return types.JID{}, fmt.Errorf("número demasiado corto: %s", phone)
	}

	jid, err := types.ParseJID(formattedPhone)
	if err != nil {
		return types.JID{}, fmt.Errorf("número inválido '%s': %v", formattedPhone, err)
	}

	log.Printf("📞 Número formateado: %s -> %s", phone, formattedPhone)
	return jid, nil
}

func (a *WhatsMeowAdapter) generateQRImage(qrCode string) (string, error) {
	log.Printf("🔍 QR string recibido: %s", qrCode) // ← VERIFICAR ESTO

	qrPNG, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("error generando QR: %v", err)
	}

	log.Printf("📊 QR PNG generado: %d bytes", len(qrPNG)) // ← VERIFICAR ESTO

	qrBase64 := base64.StdEncoding.EncodeToString(qrPNG)
	qrDataURL := "data:image/png;base64," + qrBase64

	log.Printf("🔗 Data URL generada: %d caracteres", len(qrDataURL)) // ← VERIFICAR ESTO

	return qrDataURL, nil
}

func (a *WhatsMeowAdapter) setQRCode(code string) {
	a.qrMutex.Lock()
	defer a.qrMutex.Unlock()
	a.currentQR = code
	a.isConnected = false
	log.Printf("🔐 QR Code actualizado: %s...", code[:20])
}

func (a *WhatsMeowAdapter) clearQRCode() {
	a.qrMutex.Lock()
	defer a.qrMutex.Unlock()
	a.currentQR = ""
}

func (a *WhatsMeowAdapter) setConnected(connected bool) {
	a.qrMutex.Lock()
	defer a.qrMutex.Unlock()
	a.isConnected = connected
	if connected {
		a.currentQR = ""
		log.Printf("✅ Estado cambiado a: connected")
	} else {
		log.Printf("🔌 Estado cambiado a: disconnected")
	}
}
