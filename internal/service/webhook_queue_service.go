package service

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type WebhookQueueService struct {
	repo repository.UCPWebhookQueueRepository
}

func NewWebhookQueueService(repo repository.UCPWebhookQueueRepository) *WebhookQueueService {
	return &WebhookQueueService{repo: repo}
}

func (s *WebhookQueueService) Enqueue(eventID string, payload string) error {
	if s == nil || s.repo == nil {
		return nil
	}
	job := &domain.UCPWebhookJob{
		EventID:     eventID,
		Payload:     payload,
		Status:      "pending",
		Attempts:    0,
		NextRetryAt: time.Now(),
		CreatedAt:   time.Now(),
	}
	return s.repo.Create(job)
}

func (s *WebhookQueueService) List(offset, limit int) ([]*domain.UCPWebhookJob, int64, error) {
	if s == nil || s.repo == nil {
		return []*domain.UCPWebhookJob{}, 0, nil
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

func (s *WebhookQueueService) RescheduleNow(id int64) error {
	if s == nil || s.repo == nil {
		return nil
	}
	job, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	job.NextRetryAt = time.Now()
	job.Status = "retrying"
	return s.repo.Update(job)
}
