package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type ucpWebhookReplayRepository struct {
	db *database.DB
}

func NewUCPWebhookReplayRepository(db *database.DB) UCPWebhookReplayRepository {
	return &ucpWebhookReplayRepository{db: db}
}

func (r *ucpWebhookReplayRepository) FindByHash(hash string) (*domain.UCPWebhookReplay, error) {
	var replay domain.UCPWebhookReplay
	err := r.db.Where("payload_hash = ?", hash).First(&replay).Error
	if err != nil {
		return nil, err
	}
	return &replay, nil
}

func (r *ucpWebhookReplayRepository) Create(replay *domain.UCPWebhookReplay) error {
	return r.db.Create(replay).Error
}
