package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/meowucp/internal/ucp/mock"
)

type webhookPayload struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Timestamp string `json:"timestamp"`
	Order     struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"order"`
}

func main() {
	endpoint := getEnv("WEBHOOK_ENDPOINT", "http://localhost:8081/ucp/v1/order-webhooks")
	kid := getEnv("MOCK_JWKS_KID", "mock-key")

	key, _ := mock.FixedKey(kid)

	payload := webhookPayload{
		EventID:   "evt_" + time.Now().Format("150405"),
		EventType: "order.paid",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	payload.Order.ID = "100001"
	payload.Order.Status = "paid"

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("marshal payload: %v", err)
	}

	timestamp := time.Now().Unix()
	signature, err := mock.SignPayload(key, timestamp, body)
	if err != nil {
		log.Fatalf("sign payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("UCP-Key-Id", kid)
	req.Header.Set("UCP-Signature", buildSignatureHeader(timestamp, signature))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("send webhook: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Webhook sent, status: %d", resp.StatusCode)
}

func buildSignatureHeader(timestamp int64, signature string) string {
	return "t=" + strconv.FormatInt(timestamp, 10) + ",v1=" + signature
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
