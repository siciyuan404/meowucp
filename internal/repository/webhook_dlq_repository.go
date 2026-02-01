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

func (r *webhookDLQRepository) FindByID(id int64) (*domain.WebhookDLQ, error) {
	var item domain.WebhookDLQ
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *webhookDLQRepository) List(offset, limit int) ([]*domain.WebhookDLQ, error) {
	items := []*domain.WebhookDLQ{}
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *webhookDLQRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.WebhookDLQ{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
