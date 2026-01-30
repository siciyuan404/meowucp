package config

import (
	"os"
	"testing"
)

func TestLoadUCPWebhookSkipSignatureVerify(t *testing.T) {
	content := []byte("ucp:\n  webhook:\n    skip_signature_verify: true\n")
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
	if !cfg.UCP.Webhook.SkipSignatureVerify {
		t.Fatalf("expected skip_signature_verify to be true")
	}
}
