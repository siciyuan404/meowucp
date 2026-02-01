package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type oauthClientRepository struct {
	db *database.DB
}

func NewOAuthClientRepository(db *database.DB) OAuthClientRepository {
	return &oauthClientRepository{db: db}
}

func (r *oauthClientRepository) Create(client *domain.OAuthClient) error {
	return r.db.Create(client).Error
}

func (r *oauthClientRepository) FindByClientID(clientID string) (*domain.OAuthClient, error) {
	var client domain.OAuthClient
	if err := r.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *oauthClientRepository) List(offset, limit int) ([]*domain.OAuthClient, error) {
	clients := []*domain.OAuthClient{}
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *oauthClientRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.OAuthClient{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
