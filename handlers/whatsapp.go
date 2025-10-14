package handlers

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"go.mau.fi/whatsmeow/proto/waE2E" // <-- nuevo
	"google.golang.org/protobuf/proto"
)

type ChatState struct {
	LastMessage  time.Time
	WaitingReply bool
	MessageCount int
	CurrentStep  string
}

type WhatsAppHandler struct {
	client      *whatsmeow.Client
	activeChats map[string]*ChatState
	mutex       sync.RWMutex
}

func NewWhatsAppHandler(client *whatsmeow.Client) *WhatsAppHandler {
	return &WhatsAppHandler{
		client:      client,
		activeChats: make(map[string]*ChatState),
		mutex:       sync.RWMutex{},
	}
}

func (wh *WhatsAppHandler) HandleMessage(evt *events.Message) {
	// Verificar que el mensaje sea de texto y no de grupo
	if evt.Info.Type != "text" || evt.Info.IsGroup {
		return
	}

	message := strings.TrimSpace(evt.Message.GetConversation())
	sender := evt.Info.Sender.String()

	log.Printf("💬 Mensaje de %s: %s", sender, message)

	// Limpiar chats inactivos (más de 30 minutos)
	wh.cleanInactiveChats()

	// Obtener o crear estado del chat
	state := wh.getChatState(sender)
	state.LastMessage = time.Now()
	state.MessageCount++

	// Si es el primer mensaje o reinicio
	if state.MessageCount == 1 || strings.ToLower(message) == "hola" {
		wh.sendWelcomeMessage(sender)
		state.WaitingReply = true
		state.CurrentStep = "waiting_confirmation"
		return
	}

	// Procesar respuesta según el estado actual
	switch state.CurrentStep {
	case "waiting_confirmation":
		wh.processConfirmation(sender, message, state)
	default:
		wh.sendWelcomeMessage(sender)
		state.CurrentStep = "waiting_confirmation"
		state.WaitingReply = true
	}
}

func (wh *WhatsAppHandler) sendWelcomeMessage(sender string) {
	welcomeMsg := `🎃 *¡HOLA! TE DAMOS LA BIENVENIDA* 🎃

El *31 de octubre es HALLOWEEN* y tenemos ofertas especiales para ti:

🕷️ *Disfraces exclusivos*
🍫 *Dulces y caramelos*
🎭 *Decoración terrorífica*
👻 *Kits de Halloween completos*

¿Te gustaría recibir nuestro catálogo y precios?

Responde con:
✅ *SÍ* - Para ver productos y promociones
❌ *NO* - Para salir del chat`

	wh.sendTextMessage(sender, welcomeMsg)
}

func (wh *WhatsAppHandler) processConfirmation(sender, message string, state *ChatState) {
	message = strings.ToLower(strings.TrimSpace(message))

	switch {
	case containsAny(message, []string{"sí", "si", "s", "yes", "y", "quiero", "dale", "ok", "okay"}):
		wh.sendPositiveResponse(sender)
		wh.closeChat(sender)

	case containsAny(message, []string{"no", "n", "not", "salir", "cancelar", "nah"}):
		wh.sendNegativeResponse(sender)
		wh.closeChat(sender)

	default:
		wh.sendRetryMessage(sender)
	}
}

func (wh *WhatsAppHandler) sendPositiveResponse(sender string) {
	response := `🎉 *¡EXCELENTE DECISIÓN!* 🎉

📱 *Visita nuestro catálogo online:*
https://tutienda.com/halloween

🛒 *Oferta especial HALLOWEEN:*
• 20% de descuento en disfraces
• 2x1 en dulces y caramelos
• Envío gratis en compras mayores a $50

📞 *¿Necesitas ayuda?*
Escríbenos a: ventas@tutienda.com

⏰ *Oferta válida hasta el 31 de octubre*

¡Gracias por tu interés! 👻🎃`

	wh.sendTextMessage(sender, response)
}

func (wh *WhatsAppHandler) sendNegativeResponse(sender string) {
	response := `😢 *Lamentamos que no estés interesado*

Pero recuerda que tenemos:
🎃 *Los mejores precios de Halloween*
👻 *Productos de alta calidad*
🚚 *Entrega rápida*

Si cambias de opinión, escríbenos ¡Siempre estaremos aquí!

¡Que tengas un *FELIZ HALLOWEEN*! 🎃✨`

	wh.sendTextMessage(sender, response)
}

func (wh *WhatsAppHandler) sendRetryMessage(sender string) {
	retryMsg := `🤔 *No entendí tu respuesta*

Por favor responde *SOLO* con:

✅ *SÍ* - Para ver nuestro catálogo de Halloween
❌ *NO* - Para finalizar la conversación

¿Te gustaría ver nuestros productos de Halloween?`

	wh.sendTextMessage(sender, retryMsg)
}

func (wh *WhatsAppHandler) sendTextMessage(sender, message string) {
	// Parsear JID
	jid, err := types.ParseJID(sender)
	if err != nil {
		log.Printf("❌ Error parseando JID %s: %v", sender, err)
		return
	}

	// Construir mensaje
	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	// Enviar
	_, err = wh.client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		log.Printf("❌ Error enviando mensaje a %s: %v", sender, err)
	} else {
		log.Printf("✅ Mensaje enviado a %s", sender)
	}
}

func (wh *WhatsAppHandler) getChatState(sender string) *ChatState {
	wh.mutex.Lock()
	defer wh.mutex.Unlock()

	if state, exists := wh.activeChats[sender]; exists {
		return state
	}

	state := &ChatState{
		LastMessage:  time.Now(),
		WaitingReply: false,
		MessageCount: 0,
		CurrentStep:  "initial",
	}
	wh.activeChats[sender] = state
	return state
}

func (wh *WhatsAppHandler) closeChat(sender string) {
	wh.mutex.Lock()
	defer wh.mutex.Unlock()

	if state, exists := wh.activeChats[sender]; exists {
		state.WaitingReply = false
		state.CurrentStep = "completed"
	}

	log.Printf("🔒 Chat cerrado con %s", sender)
}

func (wh *WhatsAppHandler) cleanInactiveChats() {
	wh.mutex.Lock()
	defer wh.mutex.Unlock()

	now := time.Now()
	for sender, state := range wh.activeChats {
		if now.Sub(state.LastMessage) > 30*time.Minute {
			delete(wh.activeChats, sender)
			log.Printf("🧹 Chat inactivo eliminado: %s", sender)
		}
	}
}

// Función auxiliar para verificar múltiples palabras
func containsAny(text string, words []string) bool {
	text = strings.ToLower(text)
	for _, word := range words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

// Función de registro global para el handler
func RegisterWhatsAppHandlers(client *whatsmeow.Client) {
	handler := NewWhatsAppHandler(client)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if v.Info.Type == "text" && !v.Info.IsGroup {
				go handler.HandleMessage(v)
			}
		}
	})
}
