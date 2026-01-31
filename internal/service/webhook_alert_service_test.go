package service

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeWebhookAlertRepo struct {
	items []*domain.UCPWebhookAlert
}

func (f *fakeWebhookAlertRepo) Create(alert *domain.UCPWebhookAlert) error {
	f.items = append(f.items, alert)
	return nil
}

func (f *fakeWebhookAlertRepo) List(offset, limit int) ([]*domain.UCPWebhookAlert, error) {
	if offset >= len(f.items) {
		return []*domain.UCPWebhookAlert{}, nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], nil
}

func (f *fakeWebhookAlertRepo) Count() (int64, error) {
	return int64(len(f.items)), nil
}

func (f *fakeWebhookAlertRepo) ExistsRecent(eventID, reason string, window time.Duration) (bool, error) {
	return false, nil
}

type fakeWebhookEventRepo struct {
	items map[string]*domain.UCPWebhookEvent
}

func (f *fakeWebhookEventRepo) Create(event *domain.UCPWebhookEvent) error {
	if f.items == nil {
		f.items = map[string]*domain.UCPWebhookEvent{}
	}
	f.items[event.EventID] = event
	return nil
}

func (f *fakeWebhookEventRepo) FindByEventID(eventID string) (*domain.UCPWebhookEvent, error) {
	return f.items[eventID], nil
}

func (f *fakeWebhookEventRepo) UpdateStatus(eventID string, status string) error {
	if event, ok := f.items[eventID]; ok {
		event.Status = status
	}
	return nil
}

func (f *fakeWebhookEventRepo) MarkProcessed(eventID string) error {
	if event, ok := f.items[eventID]; ok {
		now := time.Now()
		event.ProcessedAt = &now
	}
	return nil
}

func TestWebhookAlertIncludesEventContext(t *testing.T) {
	now := time.Date(2026, 1, 30, 12, 0, 0, 0, time.UTC)
	alertRepo := &fakeWebhookAlertRepo{items: []*domain.UCPWebhookAlert{
		{EventID: "evt_1", Reason: "delivery_failed", Details: "boom", Attempts: 2, CreatedAt: now},
	}}
	eventRepo := &fakeWebhookEventRepo{items: map[string]*domain.UCPWebhookEvent{
		"evt_1": {EventID: "evt_1", OrderID: "ord_1"},
	}}
	service := NewWebhookAlertService(alertRepo, eventRepo)

	alerts, _, err := service.List(0, 10)
	if err != nil {
		t.Fatalf("list alerts: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}

	var details map[string]string
	if err := json.Unmarshal([]byte(alerts[0].Details), &details); err != nil {
		t.Fatalf("expected details to be json: %v", err)
	}
	if details["event_id"] != "evt_1" {
		t.Fatalf("expected event_id in details")
	}
	if details["order_id"] != "ord_1" {
		t.Fatalf("expected order_id in details")
	}
}
