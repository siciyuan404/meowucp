package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOAuthTokenExchange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthTokenHandler()
	r := gin.New()
	r.POST("/oauth2/token", handler.Token)

	body := "grant_type=authorization_code&client_id=ucp-client&client_secret=ucp-secret&code=authcode"
	req := httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	accessToken, _ := payload["access_token"].(string)
	if accessToken == "" {
		t.Fatalf("expected access_token to be set")
	}
	if payload["token_type"] != "bearer" {
		t.Fatalf("expected token_type bearer")
	}
	expiresIn, ok := payload["expires_in"].(float64)
	if !ok || expiresIn <= 0 {
		t.Fatalf("expected expires_in to be positive")
	}
}
