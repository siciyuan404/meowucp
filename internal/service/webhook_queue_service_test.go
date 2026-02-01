package service

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeQueueRepo struct {
	jobs    map[int64]*domain.UCPWebhookJob
	updated *domain.UCPWebhookJob
}

func newFakeQueueRepo() *fakeQueueRepo {
	return &fakeQueueRepo{jobs: map[int64]*domain.UCPWebhookJob{}}
}

func (f *fakeQueueRepo) Create(job *domain.UCPWebhookJob) error {
	f.jobs[job.ID] = job
	return nil
}

func (f *fakeQueueRepo) ListDue(limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}

func (f *fakeQueueRepo) Update(job *domain.UCPWebhookJob) error {
	f.updated = job
	return nil
}

func (f *fakeQueueRepo) List(offset, limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}

func (f *fakeQueueRepo) Count() (int64, error) {
	return int64(len(f.jobs)), nil
}

func (f *fakeQueueRepo) FindByID(id int64) (*domain.UCPWebhookJob, error) {
	job, ok := f.jobs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return job, nil
}

type fakeWebhookDLQRepo struct {
	created *domain.WebhookDLQ
}

func (f *fakeWebhookDLQRepo) Create(item *domain.WebhookDLQ) error {
	f.created = item
	return nil
}

type fakeWebhookReplayLogRepo struct {
	created *domain.WebhookReplayLog
}

func (f *fakeWebhookReplayLogRepo) Create(item *domain.WebhookReplayLog) error {
	f.created = item
	return nil
}

func TestBuildOrderWebhookPayload(t *testing.T) {
	order := &domain.Order{OrderNo: "ORD-1", Status: "paid", CreatedAt: time.Unix(1000, 0).UTC()}
	payload, err := buildOrderWebhookPayload(order, "paid")
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	var decoded orderWebhookPayload
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if decoded.EventType != "paid" {
		t.Fatalf("expected event_type paid, got %s", decoded.EventType)
	}
	if decoded.OrderNo != "ORD-1" {
		t.Fatalf("expected order_no ORD-1, got %s", decoded.OrderNo)
	}
	if decoded.Status != "paid" {
		t.Fatalf("expected status paid, got %s", decoded.Status)
	}
	if decoded.CreatedAt.IsZero() {
		t.Fatalf("expected created_at to be set")
	}
}

func TestWebhookQueueRescheduleNow(t *testing.T) {
	repo := newFakeQueueRepo()
	job := &domain.UCPWebhookJob{ID: 1, EventID: "evt_1", Status: "failed", NextRetryAt: time.Now().Add(time.Hour)}
	repo.jobs[1] = job

	service := NewWebhookQueueService(repo)
	if err := service.RescheduleNow(1); err != nil {
		t.Fatalf("reschedule now: %v", err)
	}
	if repo.updated == nil {
		t.Fatalf("expected job to be updated")
	}
	if repo.updated.Status != "retrying" {
		t.Fatalf("expected status retrying, got %s", repo.updated.Status)
	}
	if !repo.updated.NextRetryAt.Before(time.Now().Add(2 * time.Second)) {
		t.Fatalf("expected next_retry_at to be near now")
	}
}

func TestWebhookQueueMovesToDLQAfterMaxAttempts(t *testing.T) {
	repo := newFakeQueueRepo()
	job := &domain.UCPWebhookJob{ID: 10, EventID: "evt_1", Status: "failed", Attempts: 5, Payload: "{}"}
	repo.jobs[10] = job

	dlqRepo := &fakeWebhookDLQRepo{}
	replayRepo := &fakeWebhookReplayLogRepo{}
	service := NewWebhookQueueServiceWithDeps(repo, dlqRepo, replayRepo)

	if err := service.MoveToDLQ(job, "max_attempts"); err != nil {
		t.Fatalf("move to dlq: %v", err)
	}
	if dlqRepo.created == nil {
		t.Fatalf("expected dlq record")
	}
}

func TestWebhookReplayLogsResult(t *testing.T) {
	repo := newFakeQueueRepo()
	job := &domain.UCPWebhookJob{ID: 11, EventID: "evt_2", Status: "failed", Attempts: 2, Payload: "{}"}
	repo.jobs[11] = job

	dlqRepo := &fakeWebhookDLQRepo{}
	replayRepo := &fakeWebhookReplayLogRepo{}
	service := NewWebhookQueueServiceWithDeps(repo, dlqRepo, replayRepo)

	if err := service.ReplayJob(11); err != nil {
		t.Fatalf("replay job: %v", err)
	}
	if replayRepo.created == nil {
		t.Fatalf("expected replay log")
	}
}
