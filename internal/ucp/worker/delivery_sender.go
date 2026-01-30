package worker

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/meowucp/internal/domain"
)

var ErrDeliveryURLMissing = errors.New("delivery_url_missing")

type DeliverySender struct {
	url    string
	client *http.Client
}

func NewDeliverySender(url string, timeout time.Duration) *DeliverySender {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &DeliverySender{
		url:    strings.TrimSpace(url),
		client: &http.Client{Timeout: timeout},
	}
}

func (s *DeliverySender) Send(job *domain.UCPWebhookJob) error {
	if s.url == "" {
		return ErrDeliveryURLMissing
	}
	if job == nil {
		return errors.New("nil_job")
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewBufferString(job.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("delivery_failed")
	}

	return nil
}
