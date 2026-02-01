package service

import (
	"errors"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type OAuthTokenService struct {
	clientRepo repository.OAuthClientRepository
	tokenRepo  repository.OAuthTokenRepository
}

func NewOAuthTokenService(clientRepo repository.OAuthClientRepository, tokenRepo repository.OAuthTokenRepository) *OAuthTokenService {
	return &OAuthTokenService{clientRepo: clientRepo, tokenRepo: tokenRepo}
}

func (s *OAuthTokenService) IssueToken(clientID, scopes string, expiresAt time.Time) (*domain.OAuthToken, error) {
	if s == nil || s.tokenRepo == nil {
		return nil, errors.New("oauth_token_repo_unavailable")
	}
	item := &domain.OAuthToken{
		Token:     clientID + ":token",
		ClientID:  clientID,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
	if err := s.tokenRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *OAuthTokenService) Revoke(token string) error {
	if s == nil || s.tokenRepo == nil {
		return errors.New("oauth_token_repo_unavailable")
	}
	return s.tokenRepo.Revoke(token, time.Now())
}
