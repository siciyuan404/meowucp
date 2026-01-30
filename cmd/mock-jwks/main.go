package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/meowucp/internal/ucp/mock"
)

func main() {
	port := getEnv("MOCK_JWKS_PORT", "9092")
	kid := getEnv("MOCK_JWKS_KID", "mock-key")

	_, jwk := mock.FixedKey(kid)
	set := mock.JWKSet{Keys: []mock.JWK{jwk}}
	encoded, err := json.Marshal(set)
	if err != nil {
		log.Fatalf("marshal jwks: %v", err)
	}

	http.HandleFunc("/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(encoded)
	})

	addr := ":" + port
	log.Printf("Mock JWKS server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start mock jwks server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
