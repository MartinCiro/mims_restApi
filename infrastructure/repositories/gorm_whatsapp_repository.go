package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"scrapper_go_email/core/whatsapp"
	"scrapper_go_email/infrastructure/database"
	"scrapper_go_email/infrastructure/database/models"

	"gorm.io/gorm"
)

type GormWhatsAppRepository struct {
	db *gorm.DB
}

func NewGormWhatsAppRepository() *GormWhatsAppRepository {
	return &GormWhatsAppRepository{
		db: database.GetDB(),
	}
}

func (r *GormWhatsAppRepository) SaveSession(ctx context.Context, session *whatsapp.Session) error {
	dbSession := models.Session{
		ID:         session.ID,
		Provider:   session.Provider,
		Status:     string(session.Status),
		QRCode:     session.QRCode,
		WebhookURL: session.WebhookURL,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
	}

	if session.ConnectedAt != nil {
		dbSession.ConnectedAt = session.ConnectedAt
	}

	// Serializar Metadata a JSON
	if session.Metadata != nil {
		metadataJSON, err := json.Marshal(session.Metadata)
		if err != nil {
			return fmt.Errorf("error serializando metadata: %v", err)
		}
		dbSession.Metadata = string(metadataJSON)
	}

	// Serializar UserInfo a JSON
	if session.UserInfo != nil {
		userInfoJSON, err := json.Marshal(session.UserInfo)
		if err != nil {
			return fmt.Errorf("error serializando user_info: %v", err)
		}
		dbSession.UserInfo = string(userInfoJSON)
	}

	return r.db.WithContext(ctx).Save(&dbSession).Error
}

func (r *GormWhatsAppRepository) FindSession(ctx context.Context, sessionID string) (*whatsapp.Session, error) {
	var dbSession models.Session
	err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&dbSession).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.convertToDomainSession(&dbSession)
}

func (r *GormWhatsAppRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Where("id = ?", sessionID).Delete(&models.Session{}).Error
}

func (r *GormWhatsAppRepository) ListSessions(ctx context.Context) ([]*whatsapp.Session, error) {
	var dbSessions []models.Session
	err := r.db.WithContext(ctx).Find(&dbSessions).Error
	if err != nil {
		return nil, err
	}

	sessions := make([]*whatsapp.Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		session, err := r.convertToDomainSession(&dbSession)
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *GormWhatsAppRepository) UpdateSessionStatus(ctx context.Context, sessionID string, status whatsapp.SessionStatus) error {
	return r.db.WithContext(ctx).Model(&models.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *GormWhatsAppRepository) convertToDomainSession(dbSession *models.Session) (*whatsapp.Session, error) {
	session := &whatsapp.Session{
		ID:          dbSession.ID,
		Provider:    dbSession.Provider,
		Status:      whatsapp.SessionStatus(dbSession.Status),
		QRCode:      dbSession.QRCode,
		WebhookURL:  dbSession.WebhookURL,
		CreatedAt:   dbSession.CreatedAt,
		UpdatedAt:   dbSession.UpdatedAt,
		ConnectedAt: dbSession.ConnectedAt,
	}

	// Deserializar Metadata
	if dbSession.Metadata != "" {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(dbSession.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("error deserializando metadata: %v", err)
		}
		session.Metadata = metadata
	}

	// Deserializar UserInfo
	if dbSession.UserInfo != "" {
		var userInfo whatsapp.UserInfo
		if err := json.Unmarshal([]byte(dbSession.UserInfo), &userInfo); err != nil {
			return nil, fmt.Errorf("error deserializando user_info: %v", err)
		}
		session.UserInfo = &userInfo
	}

	return session, nil
}
