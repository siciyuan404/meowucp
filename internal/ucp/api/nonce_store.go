package api

import (
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/security"
)

type webhookReplayNonceStore struct {
	replay *service.WebhookReplayService
}

func NewWebhookReplayNonceStore(replay *service.WebhookReplayService) security.NonceStore {
	return &webhookReplayNonceStore{replay: replay}
}

func (s *webhookReplayNonceStore) Seen(nonce string) (bool, error) {
	if s == nil || s.replay == nil {
		return false, nil
	}
	return s.replay.Seen(nonce)
}

func (s *webhookReplayNonceStore) Mark(nonce string) error {
	if s == nil || s.replay == nil {
		return nil
	}
	return s.replay.Mark(nonce, 600)
}
