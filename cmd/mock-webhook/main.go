package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/meowucp/internal/ucp/mock"
)

func main() {
	port := getEnv("MOCK_WEBHOOK_PORT", "9091")
	failStatus := parseInt(getEnv("MOCK_WEBHOOK_FAIL_STATUS", "0"))

	receiver := mock.NewReceiver(failStatus)

	http.Handle("/ucp/v1/order-webhooks", receiver)

	addr := ":" + port
	log.Printf("Mock webhook server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start mock webhook server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}
