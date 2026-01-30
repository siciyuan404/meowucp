package service

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookReplayService struct {
	repo repository.UCPWebhookReplayRepository
}

func NewWebhookReplayService(repo repository.UCPWebhookReplayRepository) *WebhookReplayService {
	return &WebhookReplayService{repo: repo}
}

func (s *WebhookReplayService) Seen(hash string) (bool, error) {
	if s == nil || s.repo == nil {
		return false, nil
	}
	_, err := s.repo.FindByHash(hash)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *WebhookReplayService) Mark(hash string, ttlSeconds int) error {
	if s == nil || s.repo == nil {
		return nil
	}
	expires := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	return s.repo.Create(&domain.UCPWebhookReplay{
		PayloadHash: hash,
		ExpiresAt:   expires,
	})
}
