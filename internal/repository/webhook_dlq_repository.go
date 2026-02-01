package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type webhookDLQRepository struct {
	db *database.DB
}

func NewWebhookDLQRepository(db *database.DB) WebhookDLQRepository {
	return &webhookDLQRepository{db: db}
}

func (r *webhookDLQRepository) Create(item *domain.WebhookDLQ) error {
	return r.db.Create(item).Error
}
