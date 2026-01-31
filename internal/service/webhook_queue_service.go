package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/model"
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

func (s *WebhookQueueService) EnqueuePaidEvent(order *domain.Order) (string, error) {
	if s == nil || s.repo == nil {
		return "", nil
	}
	if order == nil {
		return "", errors.New("order required")
	}
	eventID := "evt_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	status := order.Status
	if status == "" {
		status = "paid"
	}
	payload := model.OrderWebhookEvent{
		EventID:   eventID,
		EventType: "order.paid",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Order: model.OrderWebhookOrder{
			ID:     strconv.FormatInt(order.ID, 10),
			Status: status,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	if err := s.Enqueue(eventID, string(body)); err != nil {
		return "", err
	}
	return eventID, nil
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
