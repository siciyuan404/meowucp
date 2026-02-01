package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeAuditLogService struct {
	items []*domain.AuditLog
	total int64
}

func newFakeAuditLogService() *fakeAuditLogService {
	now := time.Now()
	return &fakeAuditLogService{
		items: []*domain.AuditLog{
			{ID: 1, Actor: "admin", Action: "update", Target: "order:1", Payload: nil, CreatedAt: now},
			{ID: 2, Actor: "admin", Action: "delete", Target: "product:2", Payload: nil, CreatedAt: now},
		},
		total: 2,
	}
}

func (f *fakeAuditLogService) List(offset, limit int) ([]*domain.AuditLog, int64, error) {
	start := offset
	if start >= len(f.items) {
		return []*domain.AuditLog{}, f.total, nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[start:end], f.total, nil
}

func TestAdminAuditList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeAuditLogService()
	handler := NewAdminAuditHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/audit-logs", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs?page=1&limit=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !containsString(resp.Body.String(), "pagination") {
		t.Fatalf("expected pagination in response")
	}
}
