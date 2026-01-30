package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWKVerifierVerifySuccess(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	keyID := "test-key"
	jwkSet := JWKSet{
		Keys: []JWK{
			{
				KTY: "EC",
				CRV: "P-256",
				X:   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.X.Bytes()),
				Y:   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.Y.Bytes()),
				KID: keyID,
				Alg: "ES256",
				Use: "sig",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(jwkSet)
	}))
	defer server.Close()

	verifier := NewJWKVerifier(server.URL, 300)
	body := []byte(`{"event_id":"evt_1"}`)
	timestamp := time.Now().Unix()
	msg := []byte(fmt.Sprintf("%d.%s", timestamp, string(body)))
	hash := sha256.Sum256(msg)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", nil)
	req.Header.Set(signatureHeader, fmt.Sprintf("t=%d,v1=%s", timestamp, base64.RawURLEncoding.EncodeToString(signature)))
	req.Header.Set(keyIDHeader, keyID)

	if err := verifier.Verify(req, body); err != nil {
		t.Fatalf("expected verification success, got %v", err)
	}
}

func TestJWKVerifierMissingSignature(t *testing.T) {
	verifier := NewJWKVerifier("https://example.com/jwks", 300)
	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", nil)

	if err := verifier.Verify(req, []byte("{}")); err == nil {
		t.Fatalf("expected error for missing signature")
	}
}
