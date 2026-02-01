package repository

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type oauthTokenRepository struct {
	db *database.DB
}

func NewOAuthTokenRepository(db *database.DB) OAuthTokenRepository {
	return &oauthTokenRepository{db: db}
}

func (r *oauthTokenRepository) Create(token *domain.OAuthToken) error {
	return r.db.Create(token).Error
}

func (r *oauthTokenRepository) FindByToken(token string) (*domain.OAuthToken, error) {
	var item domain.OAuthToken
	if err := r.db.Where("token = ?", token).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *oauthTokenRepository) Revoke(token string, revokedAt time.Time) error {
	return r.db.Model(&domain.OAuthToken{}).Where("token = ?", token).Update("revoked_at", revokedAt).Error
}
