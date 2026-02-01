package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeWebhookDLQService struct {
	items    []*domain.WebhookDLQ
	replayed int64
}

func (f *fakeWebhookDLQService) ListDLQ(offset, limit int) ([]*domain.WebhookDLQ, int64, error) {
	return f.items, int64(len(f.items)), nil
}

func (f *fakeWebhookDLQService) ReplayDLQ(id int64) error {
	f.replayed = id
	return nil
}

func TestAdminListWebhookDLQ(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeWebhookDLQService{items: []*domain.WebhookDLQ{{ID: 1, JobID: 10, Reason: "failed", Payload: "{}"}}}
	handler := NewAdminWebhookDLQHandler(service)

	r := gin.New()
	r.GET("/api/v1/admin/webhooks/dlq", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/webhooks/dlq", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}

func TestAdminReplayWebhookDLQ(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeWebhookDLQService{}
	handler := NewAdminWebhookDLQHandler(service)

	r := gin.New()
	r.POST("/api/v1/admin/webhooks/dlq/:id/replay", handler.Replay)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/webhooks/dlq/5/replay", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if service.replayed != 5 {
		t.Fatalf("expected replay to be called")
	}
}
