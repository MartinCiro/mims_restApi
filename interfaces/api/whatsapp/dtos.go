package whatsapp

// RegisterRequestDTO DTO para registro de sesión
type RegisterRequestDTO struct {
	SessionID  string            `json:"session_id" binding:"required,min=3,max=50"`
	Provider   string            `json:"provider" binding:"required,oneof=whatsmeow baileys"`
	WebhookURL string            `json:"webhook_url" binding:"omitempty,url"`
	Metadata   map[string]string `json:"metadata" binding:"omitempty"`
}

// MessageRequestDTO DTO para envío de mensaje
type MessageRequestDTO struct {
	SessionID string `json:"session_id" binding:"required"`
	To        string `json:"to" binding:"required,min=10,max=20"`
	Message   string `json:"message" binding:"required,min=1,max=1000"`
	Type      string `json:"type" binding:"required,oneof=text image document"`
}

// SessionResponseDTO DTO para respuesta de sesión
type SessionResponseDTO struct {
	ID          string            `json:"id"`
	Provider    string            `json:"provider"`
	Status      string            `json:"status"`
	QRCode      string            `json:"qr_code,omitempty"`
	UserInfo    *UserInfoDTO      `json:"user_info,omitempty"`
	WebhookURL  string            `json:"webhook_url,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	ConnectedAt string            `json:"connected_at,omitempty"`
}

// RegisterResponseDTO DTO específico para respuesta de registro
type RegisterResponseDTO struct {
	SessionID  string            `json:"session_id"`
	Provider   string            `json:"provider"`
	Status     string            `json:"status"`
	QRCode     string            `json:"qr_code,omitempty"`
	WebhookURL string            `json:"webhook_url,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  string            `json:"created_at"`
}

type UserInfoDTO struct {
	Phone      string `json:"phone"`
	Name       string `json:"name,omitempty"`
	ProfilePic string `json:"profile_pic,omitempty"`
	Status     string `json:"status,omitempty"`
}

type MessageResultDTO struct {
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	SentAt    string `json:"sent_at"`
	Status    string `json:"status"`
}
