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

type fakeOrderWebhookOrderService struct {
	order *domain.Order
}

func (f *fakeOrderWebhookOrderService) GetOrder(id int64) (*domain.Order, error) {
	return f.order, nil
}

type fakeOrderWebhookQueueService struct {
	delivered     *domain.Order
	deliveredType string
	enqueued      *domain.Order
	enqueuedType  string
}

func (f *fakeOrderWebhookQueueService) EnqueueOrderEvent(order *domain.Order, eventType string) error {
	f.enqueued = order
	f.enqueuedType = eventType
	return nil
}

func (f *fakeOrderWebhookQueueService) DeliverOrderEvent(order *domain.Order, eventType string, deliveryURL string, timeout time.Duration) error {
	f.delivered = order
	f.deliveredType = eventType
	return nil
}

func TestAdminOrderWebhookSync(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderSvc := &fakeOrderWebhookOrderService{order: &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "paid"}}
	webhookSvc := &fakeOrderWebhookQueueService{}
	handler := NewAdminOrderWebhookHandler(orderSvc, webhookSvc, AdminOrderWebhookConfig{
		DeliveryURL: "https://example.com/ucp/v1/order-webhooks",
		Timeout:     time.Second,
	})

	r := gin.New()
	r.POST("/api/v1/admin/orders/:id/webhook", handler.Trigger)

	body := `{"event_type":"paid","mode":"sync"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/1/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if webhookSvc.delivered == nil || webhookSvc.deliveredType != "paid" {
		t.Fatalf("expected DeliverOrderEvent to be called")
	}
}
