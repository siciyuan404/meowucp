package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type OAuthClientService struct {
	repo repository.OAuthClientRepository
}

func NewOAuthClientService(repo repository.OAuthClientRepository) *OAuthClientService {
	return &OAuthClientService{repo: repo}
}

func (s *OAuthClientService) Create(clientID, secret, scopes string) (*domain.OAuthClient, error) {
	if s == nil || s.repo == nil {
		return nil, errors.New("oauth_client_repo_unavailable")
	}
	if clientID == "" || secret == "" {
		return nil, errors.New("missing_client_credentials")
	}
	secretHash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	client := &domain.OAuthClient{
		ClientID:   clientID,
		SecretHash: string(secretHash),
		Scopes:     scopes,
		Status:     "active",
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Create(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *OAuthClientService) List(offset, limit int) ([]*domain.OAuthClient, int64, error) {
	if s == nil || s.repo == nil {
		return []*domain.OAuthClient{}, 0, nil
	}
	items, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
