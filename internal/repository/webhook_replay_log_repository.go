package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type webhookReplayLogRepository struct {
	db *database.DB
}

func NewWebhookReplayLogRepository(db *database.DB) WebhookReplayLogRepository {
	return &webhookReplayLogRepository{db: db}
}

func (r *webhookReplayLogRepository) Create(item *domain.WebhookReplayLog) error {
	return r.db.Create(item).Error
}
