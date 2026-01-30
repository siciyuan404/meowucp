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

type fakeOrderService struct {
	orders        map[int64]*domain.Order
	lastID        int64
	statusUpdates map[int64]string
	lastFilters   map[string]interface{}
}

func newFakeOrderService() *fakeOrderService {
	return &fakeOrderService{orders: map[int64]*domain.Order{}, statusUpdates: map[int64]string{}}
}

func (f *fakeOrderService) ListOrders(offset, limit int, filters map[string]interface{}) ([]*domain.Order, int64, error) {
	f.lastFilters = filters
	items := make([]*domain.Order, 0, len(f.orders))
	for _, item := range f.orders {
		items = append(items, item)
	}
	return items, int64(len(items)), nil
}

func (f *fakeOrderService) GetOrder(id int64) (*domain.Order, error) {
	item, ok := f.orders[id]
	if !ok {
		return nil, errNotFound
	}
	return item, nil
}

func (f *fakeOrderService) UpdateOrderStatus(id int64, status string) error {
	f.statusUpdates[id] = status
	if order, ok := f.orders[id]; ok {
		order.Status = status
	}
	return nil
}

func TestAdminOrderList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "pending", Total: 10, CreatedAt: time.Now()}
	svc.orders[2] = &domain.Order{ID: 2, OrderNo: "ORD-2", Status: "paid", Total: 20, CreatedAt: time.Now()}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/orders", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders?page=1&limit=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "pagination") {
		t.Fatalf("expected pagination in response")
	}
}

func TestAdminOrderListFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "pending", Total: 10}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/orders", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders?status=paid&order_no=ORD-1&from=2026-01-01&to=2026-01-31", nil)
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
	if _, ok := svc.lastFilters["order_no = ?"]; !ok {
		t.Fatalf("expected order_no filter")
	}
	if _, ok := svc.lastFilters["created_at >= ?"]; !ok {
		t.Fatalf("expected from filter")
	}
	if _, ok := svc.lastFilters["created_at <= ?"]; !ok {
		t.Fatalf("expected to filter")
	}
}

func TestAdminOrderListFiltersExtra(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "pending", Total: 10}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/orders", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders?user_id=10&amount_min=100&amount_max=200&sku=SKU-1", nil)
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
	if _, ok := svc.lastFilters["total >= ?"]; !ok {
		t.Fatalf("expected amount_min filter")
	}
	if _, ok := svc.lastFilters["total <= ?"]; !ok {
		t.Fatalf("expected amount_max filter")
	}
	if _, ok := svc.lastFilters["item_sku"]; !ok {
		t.Fatalf("expected sku filter")
	}
}

func TestAdminOrderGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "pending", Total: 10}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/orders/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "ORD-1") {
		t.Fatalf("expected order in response")
	}
}

func TestAdminOrderShip(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "paid"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/ship", handler.Ship)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/ship", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.statusUpdates[1] != "shipped" {
		t.Fatalf("expected status update to shipped")
	}
}

func TestAdminOrderShipRejectsInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "pending"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/ship", handler.Ship)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/ship", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}
}

func TestAdminOrderCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "pending"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/cancel", handler.Cancel)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/cancel", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.statusUpdates[1] != "cancelled" {
		t.Fatalf("expected status update to cancelled")
	}
}

func TestAdminOrderCancelRejectsInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "shipped"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/cancel", handler.Cancel)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/cancel", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}
}

func TestAdminOrderRefund(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "paid"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/refund", handler.Refund)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/refund", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.statusUpdates[1] != "refunded" {
		t.Fatalf("expected status update to refunded")
	}
}

func TestAdminOrderRefundRejectsInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeOrderService()
	svc.orders[1] = &domain.Order{ID: 1, Status: "pending"}

	handler := NewAdminOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/refund", handler.Refund)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/refund", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}
}
