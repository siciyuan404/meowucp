package config

import (
	"os"
	"testing"
)

func TestLoadUCPWebhookAlertPolicy(t *testing.T) {
	content := []byte("ucp:\n  webhook:\n    alert_min_attempts: 3\n    alert_dedupe_seconds: 120\n")
	file, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(content); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_ = file.Close()

	cfg, err := Load(file.Name())
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.UCP.Webhook.AlertMinAttempts != 3 {
		t.Fatalf("expected alert_min_attempts 3")
	}
	if cfg.UCP.Webhook.AlertDedupeSeconds != 120 {
		t.Fatalf("expected alert_dedupe_seconds 120")
	}
}
