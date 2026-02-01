package service

import (
	"encoding/json"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookAlertService struct {
	repo      repository.UCPWebhookAlertRepository
	eventRepo repository.UCPWebhookEventRepository
}

func NewWebhookAlertService(repo repository.UCPWebhookAlertRepository, eventRepo repository.UCPWebhookEventRepository) *WebhookAlertService {
	return &WebhookAlertService{repo: repo, eventRepo: eventRepo}
}

func (s *WebhookAlertService) Create(alert *domain.UCPWebhookAlert) error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Create(alert)
}

func (s *WebhookAlertService) CreateDedup(alert *domain.UCPWebhookAlert, window time.Duration) error {
	if s == nil || s.repo == nil {
		return nil
	}
	if alert == nil {
		return nil
	}
	if window > 0 {
		exists, err := s.repo.ExistsRecent(alert.EventID, alert.Reason, window)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
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
	for _, alert := range items {
		if alert == nil {
			continue
		}
		details := map[string]string{
			"event_id": alert.EventID,
		}
		if alert.Details != "" {
			details["error"] = alert.Details
		}
		if s.eventRepo != nil && alert.EventID != "" {
			if event, lookupErr := s.eventRepo.FindByEventID(alert.EventID); lookupErr == nil && event != nil {
				details["order_id"] = event.OrderID
			}
		}
		if payload, marshalErr := json.Marshal(details); marshalErr == nil {
			alert.Details = string(payload)
		}
	}
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
