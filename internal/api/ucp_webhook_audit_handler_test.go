package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeWebhookAuditLister struct {
	items []*domain.UCPWebhookAudit
}

func (f *fakeWebhookAuditLister) List(offset, limit int) ([]*domain.UCPWebhookAudit, int64, error) {
	if offset >= len(f.items) {
		return []*domain.UCPWebhookAudit{}, int64(len(f.items)), nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], int64(len(f.items)), nil
}

func TestListWebhookAudits(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := []*domain.UCPWebhookAudit{
		{EventID: "evt_1", Reason: "invalid_signature"},
		{EventID: "evt_2", Reason: "invalid_signature"},
	}
	lister := &fakeWebhookAuditLister{items: items}
	handler := NewWebhookAuditHandler(lister)

	r := gin.New()
	r.GET("/api/v1/admin/ucp/webhook-audits", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ucp/webhook-audits?page=2&limit=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	body := resp.Body.String()
	if !containsAll(body, []string{"evt_2", "pagination", "total"}) {
		t.Fatalf("expected response to include audit and pagination")
	}
}

func containsAll(body string, parts []string) bool {
	for _, part := range parts {
		if !strings.Contains(body, part) {
			return false
		}
	}
	return true
}
