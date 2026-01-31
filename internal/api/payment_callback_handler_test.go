package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakePaymentCallbackPaymentService struct {
	called        bool
	orderID       int64
	transactionID string
}

func (f *fakePaymentCallbackPaymentService) MarkPaymentPaid(orderID int64, transactionID string) error {
	f.called = true
	f.orderID = orderID
	f.transactionID = transactionID
	return nil
}

type fakePaymentCallbackOrderService struct {
	called bool
	id     int64
	status string
}

func (f *fakePaymentCallbackOrderService) UpdateOrderStatus(id int64, status string) error {
	f.called = true
	f.id = id
	f.status = status
	return nil
}

func TestPaymentCallbackMarksOrderPaid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentService := &fakePaymentCallbackPaymentService{}
	orderService := &fakePaymentCallbackOrderService{}
	handler := NewPaymentCallbackHandler(paymentService, orderService)

	r := gin.New()
	r.POST("/api/v1/payment/callback", handler.Handle)

	body := `{"order_id": 123, "transaction_id": "tx_001"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/callback", strings.NewReader(body))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !paymentService.called || paymentService.orderID != 123 || paymentService.transactionID != "tx_001" {
		t.Fatalf("expected payment to be marked paid")
	}
	if !orderService.called || orderService.id != 123 || orderService.status != "paid" {
		t.Fatalf("expected order status to be updated to paid")
	}
}
