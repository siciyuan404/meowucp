package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookAlertService struct {
	repo repository.UCPWebhookAlertRepository
}

func NewWebhookAlertService(repo repository.UCPWebhookAlertRepository) *WebhookAlertService {
	return &WebhookAlertService{repo: repo}
}

func (s *WebhookAlertService) Create(alert *domain.UCPWebhookAlert) error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Create(alert)
}

func (s *WebhookAlertService) List(offset, limit int) ([]*domain.UCPWebhookAlert, int64, error) {
	if s == nil || s.repo == nil {
		return []*domain.UCPWebhookAlert{}, 0, nil
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
