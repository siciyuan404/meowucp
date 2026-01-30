package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeWebhookJobAdmin struct {
	jobs        []*domain.UCPWebhookJob
	rescheduled string
}

func (f *fakeWebhookJobAdmin) List(offset, limit int) ([]*domain.UCPWebhookJob, int64, error) {
	if offset >= len(f.jobs) {
		return []*domain.UCPWebhookJob{}, int64(len(f.jobs)), nil
	}
	end := offset + limit
	if end > len(f.jobs) {
		end = len(f.jobs)
	}
	return f.jobs[offset:end], int64(len(f.jobs)), nil
}

func (f *fakeWebhookJobAdmin) RescheduleNow(id int64) error {
	for _, job := range f.jobs {
		if job.ID == id {
			job.NextRetryAt = time.Now()
			job.Status = "retrying"
			f.rescheduled = job.EventID
			return nil
		}
	}
	return nil
}

func TestListWebhookJobs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := []*domain.UCPWebhookJob{
		{ID: 1, EventID: "evt_1", Status: "retrying"},
		{ID: 2, EventID: "evt_2", Status: "failed"},
	}
	admin := &fakeWebhookJobAdmin{jobs: items}
	handler := NewWebhookJobHandler(admin)

	r := gin.New()
	r.GET("/api/v1/admin/ucp/webhook-jobs", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ucp/webhook-jobs?page=2&limit=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	body := resp.Body.String()
	if !containsAll(body, []string{"evt_2", "pagination", "total"}) {
		t.Fatalf("expected response to include job and pagination")
	}
}

func TestRescheduleWebhookJob(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := []*domain.UCPWebhookJob{
		{ID: 1, EventID: "evt_1", Status: "retrying"},
	}
	admin := &fakeWebhookJobAdmin{jobs: items}
	handler := NewWebhookJobHandler(admin)

	r := gin.New()
	r.POST("/api/v1/admin/ucp/webhook-jobs/:id/retry", handler.Retry)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/ucp/webhook-jobs/1/retry", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if admin.rescheduled != "evt_1" {
		t.Fatalf("expected reschedule to be called")
	}
}
