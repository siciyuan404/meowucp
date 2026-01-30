package repository

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type ucpWebhookEventRepository struct {
	db *database.DB
}

func NewUCPWebhookEventRepository(db *database.DB) UCPWebhookEventRepository {
	return &ucpWebhookEventRepository{db: db}
}

func (r *ucpWebhookEventRepository) Create(event *domain.UCPWebhookEvent) error {
	return r.db.Create(event).Error
}

func (r *ucpWebhookEventRepository) FindByEventID(eventID string) (*domain.UCPWebhookEvent, error) {
	var event domain.UCPWebhookEvent
	err := r.db.Where("event_id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *ucpWebhookEventRepository) UpdateStatus(eventID string, status string) error {
	return r.db.Model(&domain.UCPWebhookEvent{}).
		Where("event_id = ?", eventID).
		Update("status", status).Error
}

func (r *ucpWebhookEventRepository) MarkProcessed(eventID string) error {
	return r.db.Model(&domain.UCPWebhookEvent{}).
		Where("event_id = ?", eventID).
		Updates(map[string]interface{}{"status": "processed", "processed_at": time.Now()}).Error
}
