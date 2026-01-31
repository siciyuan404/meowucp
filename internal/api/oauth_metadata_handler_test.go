package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOAuthMetadataWellKnown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthMetadataHandler()
	r := gin.New()
	r.GET("/.well-known/oauth-authorization-server", handler.WellKnown)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/oauth-authorization-server", nil)
	req.Host = "example.com"
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	issuer, _ := payload["issuer"].(string)
	if issuer != "http://example.com" {
		t.Fatalf("expected issuer to be http://example.com")
	}
	if payload["token_endpoint"] != "http://example.com/oauth2/token" {
		t.Fatalf("expected token_endpoint to be set")
	}
	if payload["authorization_endpoint"] != "http://example.com/oauth2/authorize" {
		t.Fatalf("expected authorization_endpoint to be set")
	}

	scopes, ok := payload["scopes_supported"].([]interface{})
	if !ok || len(scopes) == 0 {
		t.Fatalf("expected scopes_supported to be present")
	}
	found := false
	for _, scope := range scopes {
		if scope == "ucp:scopes:checkout_session" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected ucp:scopes:checkout_session in scopes_supported")
	}
}
