package worker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

func TestDeliverySenderSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sender := NewDeliverySender(server.URL, 2*time.Second)
	job := &domain.UCPWebhookJob{EventID: "evt_1", Payload: "{}"}

	if err := sender.Send(job); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestDeliverySenderFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sender := NewDeliverySender(server.URL, 2*time.Second)
	job := &domain.UCPWebhookJob{EventID: "evt_1", Payload: "{}"}

	if err := sender.Send(job); err == nil {
		t.Fatalf("expected error on non-2xx")
	}
}

func TestDeliverySenderRejectsEmptyURL(t *testing.T) {
	sender := NewDeliverySender("", 2*time.Second)
	job := &domain.UCPWebhookJob{EventID: "evt_1", Payload: "{}"}

	if err := sender.Send(job); !errors.Is(err, ErrDeliveryURLMissing) {
		t.Fatalf("expected ErrDeliveryURLMissing")
	}
}
