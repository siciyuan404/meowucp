package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
)

type fakeOrderCreator struct {
	lastUserID        int64
	lastIdempotency   string
	lastShipping      string
	lastBilling       string
	lastNotes         string
	lastPaymentMethod string
	order             *domain.Order
	err               error
}

func (f *fakeOrderCreator) CreateOrder(userID int64, idempotencyKey string, shippingAddress, billingAddress, notes string, paymentMethod string) (*domain.Order, error) {
	f.lastUserID = userID
	f.lastIdempotency = idempotencyKey
	f.lastShipping = shippingAddress
	f.lastBilling = billingAddress
	f.lastNotes = notes
	f.lastPaymentMethod = paymentMethod
	if f.err != nil {
		return nil, f.err
	}
	return f.order, nil
}

func TestOrderCreateSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeOrderCreator{order: &domain.Order{ID: 10, OrderNo: "ORD-10"}}
	handler := NewOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/orders", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{
  "user_id": 10,
  "shipping_address": "Ship",
  "billing_address": "Bill",
  "payment_method": "card",
  "notes": "hi"
}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "key-123")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "ORD-10") {
		t.Fatalf("expected order in response")
	}
	if svc.lastIdempotency != "key-123" {
		t.Fatalf("expected idempotency key to be passed")
	}
}

func TestOrderCreateMissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeOrderCreator{order: &domain.Order{ID: 10, OrderNo: "ORD-10"}}
	handler := NewOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/orders", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{
  "shipping_address": "Ship",
  "billing_address": "Bill",
  "payment_method": "card"
}`))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}
}

func TestOrderCreateIdempotencyConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeOrderCreator{err: service.ErrOrderIdempotencyConflict}
	handler := NewOrderHandler(svc)

	r := gin.New()
	r.POST("/api/v1/orders", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{
  "user_id": 12,
  "shipping_address": "Ship",
  "billing_address": "Bill",
  "payment_method": "card"
}`))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "idempotency_conflict") {
		t.Fatalf("expected idempotency conflict code")
	}
}
