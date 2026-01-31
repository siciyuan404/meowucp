package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/worker"
)

type WebhookQueueService struct {
	repo repository.UCPWebhookQueueRepository
}

type orderWebhookPayload struct {
	EventType string    `json:"event_type"`
	OrderNo   string    `json:"order_no"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
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

func (s *WebhookQueueService) EnqueueOrderEvent(order *domain.Order, eventType string) error {
	if s == nil || s.repo == nil {
		return nil
	}
	payload, err := buildOrderWebhookPayload(order, eventType)
	if err != nil {
		return err
	}
	return s.Enqueue(buildOrderWebhookEventID(order, eventType), string(payload))
}

func (s *WebhookQueueService) DeliverOrderEvent(order *domain.Order, eventType string, deliveryURL string, timeout time.Duration) error {
	payload, err := buildOrderWebhookPayload(order, eventType)
	if err != nil {
		return err
	}
	job := &domain.UCPWebhookJob{
		EventID: buildOrderWebhookEventID(order, eventType),
		Payload: string(payload),
	}
	sender := worker.NewDeliverySender(deliveryURL, timeout)
	return sender.Send(job)
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

func buildOrderWebhookPayload(order *domain.Order, eventType string) ([]byte, error) {
	if order == nil {
		return nil, errors.New("order_missing")
	}
	createdAt := order.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	payload := orderWebhookPayload{
		EventType: eventType,
		OrderNo:   order.OrderNo,
		Status:    order.Status,
		CreatedAt: createdAt,
	}
	return json.Marshal(payload)
}

func buildOrderWebhookEventID(order *domain.Order, eventType string) string {
	if order == nil {
		return fmt.Sprintf("order_unknown_%s_%d", eventType, time.Now().UnixNano())
	}
	return fmt.Sprintf("order_%s_%s_%d", order.OrderNo, eventType, time.Now().UnixNano())
}
