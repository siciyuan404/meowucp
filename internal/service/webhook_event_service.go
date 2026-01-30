package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookEventService struct {
	repo repository.UCPWebhookEventRepository
}

func NewWebhookEventService(repo repository.UCPWebhookEventRepository) *WebhookEventService {
	return &WebhookEventService{repo: repo}
}

func (s *WebhookEventService) Create(event *domain.UCPWebhookEvent) error {
	return s.repo.Create(event)
}

func (s *WebhookEventService) GetByEventID(eventID string) (*domain.UCPWebhookEvent, error) {
	return s.repo.FindByEventID(eventID)
}

func (s *WebhookEventService) UpdateStatus(eventID string, status string) error {
	return s.repo.UpdateStatus(eventID, status)
}

func (s *WebhookEventService) MarkProcessed(eventID string) error {
	return s.repo.MarkProcessed(eventID)
}
