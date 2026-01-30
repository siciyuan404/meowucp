package mock

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestSignPayloadProducesVerifiableSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	payload := []byte(`{"event_id":"evt_1"}`)
	timestamp := int64(1738134123)

	sig, err := SignPayload(privateKey, timestamp, payload)
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(sig)
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}

	message := buildSignatureMessage(timestamp, payload)
	hash := sha256.Sum256(message)
	if !ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], decoded) {
		t.Fatalf("expected signature to verify")
	}
}
