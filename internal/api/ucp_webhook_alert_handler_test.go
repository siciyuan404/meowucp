package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeWebhookAlertLister struct {
	items []*domain.UCPWebhookAlert
}

func (f *fakeWebhookAlertLister) List(offset, limit int) ([]*domain.UCPWebhookAlert, int64, error) {
	if offset >= len(f.items) {
		return []*domain.UCPWebhookAlert{}, int64(len(f.items)), nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], int64(len(f.items)), nil
}

func TestListWebhookAlerts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := []*domain.UCPWebhookAlert{
		{EventID: "evt_1", Reason: "delivery_failed"},
		{EventID: "evt_2", Reason: "delivery_failed"},
	}
	lister := &fakeWebhookAlertLister{items: items}
	handler := NewWebhookAlertHandler(lister)

	r := gin.New()
	r.GET("/api/v1/admin/ucp/webhook-alerts", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ucp/webhook-alerts?page=2&limit=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	body := resp.Body.String()
	if !containsAll(body, []string{"evt_2", "pagination", "total"}) {
		t.Fatalf("expected response to include alert and pagination")
	}
}
