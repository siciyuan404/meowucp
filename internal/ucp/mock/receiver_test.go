package mock

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReceiverSuccess(t *testing.T) {
	receiver := NewReceiver(0)

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", nil)
	resp := httptest.NewRecorder()
	receiver.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}

func TestReceiverFailureStatus(t *testing.T) {
	receiver := NewReceiver(500)

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", nil)
	resp := httptest.NewRecorder()
	receiver.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.Code)
	}
}
