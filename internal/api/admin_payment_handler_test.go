package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakePaymentService struct {
	items       []*domain.Payment
	lastFilters map[string]interface{}
}

func (f *fakePaymentService) ListPayments(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, int64, error) {
	f.lastFilters = filters
	if offset >= len(f.items) {
		return []*domain.Payment{}, int64(len(f.items)), nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], int64(len(f.items)), nil
}

func TestAdminPaymentList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakePaymentService{items: []*domain.Payment{
		{ID: 1, OrderID: 10, PaymentMethod: "nowpayments", Status: "paid", Amount: 10, CreatedAt: time.Now()},
		{ID: 2, OrderID: 11, PaymentMethod: "nowpayments", Status: "failed", Amount: 20, CreatedAt: time.Now()},
	}}

	handler := NewAdminPaymentHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/payments", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/payments?page=1&limit=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "pagination") {
		t.Fatalf("expected pagination in response")
	}
}

func TestAdminPaymentListFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakePaymentService{items: []*domain.Payment{{ID: 1}}}
	handler := NewAdminPaymentHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/payments", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/payments?status=paid&method=nowpayments&from=2026-01-01&to=2026-01-31", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.lastFilters == nil {
		t.Fatalf("expected filters to be passed")
	}
	if _, ok := svc.lastFilters["status = ?"]; !ok {
		t.Fatalf("expected status filter")
	}
	if _, ok := svc.lastFilters["payment_method = ?"]; !ok {
		t.Fatalf("expected method filter")
	}
	if _, ok := svc.lastFilters["created_at >= ?"]; !ok {
		t.Fatalf("expected from filter")
	}
	if _, ok := svc.lastFilters["created_at <= ?"]; !ok {
		t.Fatalf("expected to filter")
	}
}

func TestAdminPaymentListFiltersOrderAndTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakePaymentService{items: []*domain.Payment{{ID: 1}}}
	handler := NewAdminPaymentHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/payments", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/payments?order_id=123&transaction_id=tx_001", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.lastFilters == nil {
		t.Fatalf("expected filters to be passed")
	}
	if _, ok := svc.lastFilters["order_id = ?"]; !ok {
		t.Fatalf("expected order_id filter")
	}
	if _, ok := svc.lastFilters["transaction_id = ?"]; !ok {
		t.Fatalf("expected transaction_id filter")
	}
}

func TestAdminPaymentListFiltersExtra(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakePaymentService{items: []*domain.Payment{{ID: 1}}}
	handler := NewAdminPaymentHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/payments", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/payments?user_id=10&amount_min=1&amount_max=9", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.lastFilters == nil {
		t.Fatalf("expected filters to be passed")
	}
	if _, ok := svc.lastFilters["user_id = ?"]; !ok {
		t.Fatalf("expected user_id filter")
	}
	if _, ok := svc.lastFilters["amount >= ?"]; !ok {
		t.Fatalf("expected amount_min filter")
	}
	if _, ok := svc.lastFilters["amount <= ?"]; !ok {
		t.Fatalf("expected amount_max filter")
	}
}

func TestAdminPaymentListFiltersCurrencyAndAmountPrecision(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakePaymentService{items: []*domain.Payment{{ID: 1}}}
	handler := NewAdminPaymentHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/payments", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/payments?currency=CNY&amount_min=1.005&amount_max=2.004", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.lastFilters == nil {
		t.Fatalf("expected filters to be passed")
	}
	if _, ok := svc.lastFilters["currency = ?"]; !ok {
		t.Fatalf("expected currency filter")
	}
	if value, ok := svc.lastFilters["amount >= ?"].(float64); !ok || value != 1.01 {
		t.Fatalf("expected amount_min rounded to 1.01")
	}
	if value, ok := svc.lastFilters["amount <= ?"].(float64); !ok || value != 2.0 {
		t.Fatalf("expected amount_max rounded to 2.0")
	}
}
