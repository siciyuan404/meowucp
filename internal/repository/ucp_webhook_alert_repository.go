package repository

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type ucpWebhookAlertRepository struct {
	db *database.DB
}

func NewUCPWebhookAlertRepository(db *database.DB) UCPWebhookAlertRepository {
	return &ucpWebhookAlertRepository{db: db}
}

func (r *ucpWebhookAlertRepository) Create(alert *domain.UCPWebhookAlert) error {
	return r.db.Create(alert).Error
}

func (r *ucpWebhookAlertRepository) List(offset, limit int) ([]*domain.UCPWebhookAlert, error) {
	var alerts []*domain.UCPWebhookAlert
	err := r.db.Offset(offset).Limit(limit).Order("id desc").Find(&alerts).Error
	return alerts, err
}

func (r *ucpWebhookAlertRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.UCPWebhookAlert{}).Count(&count).Error
	return count, err
}

func (r *ucpWebhookAlertRepository) ExistsRecent(eventID, reason string, window time.Duration) (bool, error) {
	var count int64
	cutoff := time.Now().Add(-window)
	query := r.db.Model(&domain.UCPWebhookAlert{}).Where("created_at >= ?", cutoff)
	if eventID != "" {
		query = query.Where("event_id = ?", eventID)
	}
	if reason != "" {
		query = query.Where("reason = ?", reason)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
