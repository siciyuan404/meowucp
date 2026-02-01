package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/model"
	"github.com/meowucp/internal/ucp/worker"
)

type WebhookQueueService struct {
	repo          repository.UCPWebhookQueueRepository
	dlqRepo       repository.WebhookDLQRepository
	replayLogRepo repository.WebhookReplayLogRepository
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

func NewWebhookQueueServiceWithDeps(repo repository.UCPWebhookQueueRepository, dlqRepo repository.WebhookDLQRepository, replayLogRepo repository.WebhookReplayLogRepository) *WebhookQueueService {
	return &WebhookQueueService{repo: repo, dlqRepo: dlqRepo, replayLogRepo: replayLogRepo}
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

func (s *WebhookQueueService) MoveToDLQ(job *domain.UCPWebhookJob, reason string) error {
	if s == nil || s.dlqRepo == nil || job == nil {
		return nil
	}
	return s.dlqRepo.Create(&domain.WebhookDLQ{
		JobID:     job.ID,
		Reason:    reason,
		Payload:   job.Payload,
		CreatedAt: time.Now(),
	})
}

func (s *WebhookQueueService) ReplayJob(jobID int64) error {
	if s == nil || s.repo == nil || s.replayLogRepo == nil {
		return nil
	}
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		return err
	}
	job.Status = "retrying"
	job.NextRetryAt = time.Now()
	if err := s.repo.Update(job); err != nil {
		return err
	}
	return s.replayLogRepo.Create(&domain.WebhookReplayLog{
		JobID:    job.ID,
		ReplayAt: time.Now(),
		Result:   "scheduled",
	})
}

type WebhookDLQService struct {
	queue *WebhookQueueService
	repo  repository.WebhookDLQRepository
}

func NewWebhookDLQService(queue *WebhookQueueService, repo repository.WebhookDLQRepository) *WebhookDLQService {
	return &WebhookDLQService{queue: queue, repo: repo}
}

func (s *WebhookDLQService) ListDLQ(offset, limit int) ([]*domain.WebhookDLQ, int64, error) {
	if s == nil || s.repo == nil {
		return []*domain.WebhookDLQ{}, 0, nil
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

func (s *WebhookDLQService) ReplayDLQ(id int64) error {
	if s == nil || s.queue == nil || s.repo == nil {
		return nil
	}
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.queue.ReplayJob(item.JobID)
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
