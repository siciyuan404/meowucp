package worker

import (
	"errors"
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeQueueStore struct {
	jobs        []*domain.UCPWebhookJob
	updatedJobs []*domain.UCPWebhookJob
}

func (f *fakeQueueStore) ListDue(limit int) ([]*domain.UCPWebhookJob, error) {
	if len(f.jobs) > limit {
		return f.jobs[:limit], nil
	}
	return f.jobs, nil
}

func (f *fakeQueueStore) Update(job *domain.UCPWebhookJob) error {
	f.updatedJobs = append(f.updatedJobs, job)
	return nil
}

type fakeAlertSink struct {
	count int
	last  *domain.UCPWebhookAlert
}

func (f *fakeAlertSink) Notify(alert *domain.UCPWebhookAlert) error {
	f.count++
	f.last = alert
	return nil
}

func TestProcessorMarksSuccess(t *testing.T) {
	now := time.Date(2026, 1, 29, 10, 0, 0, 0, time.UTC)
	store := &fakeQueueStore{
		jobs: []*domain.UCPWebhookJob{
			{ID: 1, EventID: "evt_1", Status: "pending", Attempts: 0, NextRetryAt: now},
		},
	}
	processor := NewProcessor(store, ProcessorConfig{MaxAttempts: 3})
	processor.now = func() time.Time { return now }

	processed, err := processor.ProcessOnce(func(job *domain.UCPWebhookJob) error {
		return nil
	})
	if err != nil {
		t.Fatalf("process once: %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed, got %d", processed)
	}
	if len(store.updatedJobs) != 1 {
		t.Fatalf("expected job to be updated")
	}
	updated := store.updatedJobs[0]
	if updated.Status != "processed" {
		t.Fatalf("expected status processed, got %s", updated.Status)
	}
	if updated.Attempts != 1 {
		t.Fatalf("expected attempts 1, got %d", updated.Attempts)
	}
}

func TestProcessorRetriesOnFailure(t *testing.T) {
	now := time.Date(2026, 1, 29, 10, 0, 0, 0, time.UTC)
	store := &fakeQueueStore{
		jobs: []*domain.UCPWebhookJob{
			{ID: 1, EventID: "evt_1", Status: "pending", Attempts: 0, NextRetryAt: now},
		},
	}
	processor := NewProcessor(store, ProcessorConfig{MaxAttempts: 3})
	processor.now = func() time.Time { return now }

	_, err := processor.ProcessOnce(func(job *domain.UCPWebhookJob) error {
		return errors.New("boom")
	})
	if err != nil {
		t.Fatalf("process once: %v", err)
	}

	updated := store.updatedJobs[0]
	if updated.Status != "retrying" {
		t.Fatalf("expected status retrying, got %s", updated.Status)
	}
	if updated.Attempts != 1 {
		t.Fatalf("expected attempts 1, got %d", updated.Attempts)
	}
	if !updated.NextRetryAt.After(now) {
		t.Fatalf("expected next_retry_at to be in the future")
	}
	if updated.LastError != "boom" {
		t.Fatalf("expected last_error to be set")
	}
	if updated.LastAttemptAt.IsZero() {
		t.Fatalf("expected last_attempt_at to be set")
	}
}

func TestProcessorFailsAfterMaxAttempts(t *testing.T) {
	now := time.Date(2026, 1, 29, 10, 0, 0, 0, time.UTC)
	store := &fakeQueueStore{
		jobs: []*domain.UCPWebhookJob{
			{ID: 1, EventID: "evt_1", Status: "retrying", Attempts: 2, NextRetryAt: now},
		},
	}
	processor := NewProcessor(store, ProcessorConfig{MaxAttempts: 3})
	processor.now = func() time.Time { return now }

	_, err := processor.ProcessOnce(func(job *domain.UCPWebhookJob) error {
		return errors.New("boom")
	})
	if err != nil {
		t.Fatalf("process once: %v", err)
	}

	updated := store.updatedJobs[0]
	if updated.Status != "failed" {
		t.Fatalf("expected status failed, got %s", updated.Status)
	}
	if updated.Attempts != 3 {
		t.Fatalf("expected attempts 3, got %d", updated.Attempts)
	}
}

func TestProcessorNotifiesAlertOnFailure(t *testing.T) {
	now := time.Date(2026, 1, 29, 10, 0, 0, 0, time.UTC)
	store := &fakeQueueStore{
		jobs: []*domain.UCPWebhookJob{
			{ID: 1, EventID: "evt_1", Status: "pending", Attempts: 0, NextRetryAt: now},
		},
	}
	alertSink := &fakeAlertSink{}
	processor := NewProcessor(store, ProcessorConfig{MaxAttempts: 3})
	processor.now = func() time.Time { return now }
	processor.alertSink = alertSink

	_, err := processor.ProcessOnce(func(job *domain.UCPWebhookJob) error {
		return errors.New("boom")
	})
	if err != nil {
		t.Fatalf("process once: %v", err)
	}

	if alertSink.count != 1 {
		t.Fatalf("expected alert to be sent")
	}
	if alertSink.last == nil || alertSink.last.EventID != "evt_1" {
		t.Fatalf("expected alert to include event_id")
	}
}
