package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookAuditService struct {
	repo repository.UCPWebhookAuditRepository
}

func NewWebhookAuditService(repo repository.UCPWebhookAuditRepository) *WebhookAuditService {
	return &WebhookAuditService{repo: repo}
}

func (s *WebhookAuditService) Create(audit *domain.UCPWebhookAudit) error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Create(audit)
}

func (s *WebhookAuditService) List(offset, limit int) ([]*domain.UCPWebhookAudit, int64, error) {
	if s == nil || s.repo == nil {
		return []*domain.UCPWebhookAudit{}, 0, nil
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
