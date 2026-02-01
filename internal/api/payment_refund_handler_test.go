package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeRefundPaymentService struct {
	refund *domain.PaymentRefund
}

func (f *fakeRefundPaymentService) CreateRefund(paymentID int64, amount float64, reason string) (*domain.PaymentRefund, error) {
	f.refund = &domain.PaymentRefund{PaymentID: paymentID, Amount: amount, Reason: reason}
	return f.refund, nil
}

func TestPaymentRefundHandlerCreatesRefund(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeRefundPaymentService{}
	handler := NewPaymentRefundHandler(service)

	r := gin.New()
	r.POST("/api/v1/payments/:id/refund", handler.Create)

	body := `{"amount":40,"reason":"customer_request"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/10/refund", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if service.refund == nil || service.refund.PaymentID != 10 {
		t.Fatalf("expected refund to be created")
	}
}
