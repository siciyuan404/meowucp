package worker

import (
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeAlertRepo struct {
	count  int
	last   *domain.UCPWebhookAlert
	recent bool
}

func (f *fakeAlertRepo) Create(alert *domain.UCPWebhookAlert) error {
	f.count++
	f.last = alert
	return nil
}

func (f *fakeAlertRepo) ExistsRecent(eventID, reason string, window time.Duration) (bool, error) {
	return f.recent, nil
}

func (f *fakeAlertRepo) List(offset, limit int) ([]*domain.UCPWebhookAlert, error) {
	return []*domain.UCPWebhookAlert{}, nil
}

func (f *fakeAlertRepo) Count() (int64, error) {
	return int64(f.count), nil
}

func TestDBAlertSinkCreatesAlert(t *testing.T) {
	repo := &fakeAlertRepo{}
	sink := NewAlertPolicySink(repo, AlertPolicy{MinAttempts: 1, DedupeWindow: time.Minute})
	alert := &domain.UCPWebhookAlert{
		EventID:   "evt_1",
		Reason:    "delivery_failed",
		Details:   "boom",
		Attempts:  1,
		CreatedAt: time.Now(),
	}

	if err := sink.Notify(alert); err != nil {
		t.Fatalf("notify: %v", err)
	}
	if repo.count != 1 {
		t.Fatalf("expected alert to be stored")
	}
	if repo.last.EventID != "evt_1" {
		t.Fatalf("expected event_id to be stored")
	}
}

func TestAlertSinkSkipsWithinDedupeWindow(t *testing.T) {
	repo := &fakeAlertRepo{recent: true}
	sink := NewAlertPolicySink(repo, AlertPolicy{MinAttempts: 1, DedupeWindow: time.Minute})
	alert := &domain.UCPWebhookAlert{
		EventID:   "evt_1",
		Reason:    "delivery_failed",
		Details:   "boom",
		Attempts:  1,
		CreatedAt: time.Now(),
	}

	if err := sink.Notify(alert); err != nil {
		t.Fatalf("notify: %v", err)
	}
	if repo.count != 0 {
		t.Fatalf("expected no alert when dedupe hits")
	}
}

func TestAlertSinkSkipsBelowThreshold(t *testing.T) {
	repo := &fakeAlertRepo{}
	sink := NewAlertPolicySink(repo, AlertPolicy{MinAttempts: 3, DedupeWindow: time.Minute})
	alert := &domain.UCPWebhookAlert{
		EventID:   "evt_1",
		Reason:    "delivery_failed",
		Details:   "boom",
		Attempts:  1,
		CreatedAt: time.Now(),
	}

	if err := sink.Notify(alert); err != nil {
		t.Fatalf("notify: %v", err)
	}
	if repo.count != 0 {
		t.Fatalf("expected no alert when below threshold")
	}
}
