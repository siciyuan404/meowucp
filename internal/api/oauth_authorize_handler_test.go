package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOAuthAuthorizeRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthAuthorizeHandler()
	r := gin.New()
	r.GET("/oauth2/authorize", handler.Authorize)

	url := "/oauth2/authorize?response_type=code&client_id=ucp-client&redirect_uri=https%3A%2F%2Fexample.com%2Fcb&state=abc"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", resp.Code)
	}
	location := resp.Header().Get("Location")
	if !strings.HasPrefix(location, "https://example.com/cb") {
		t.Fatalf("expected redirect to callback")
	}
	if !strings.Contains(location, "code=") {
		t.Fatalf("expected code in redirect")
	}
	if !strings.Contains(location, "state=abc") {
		t.Fatalf("expected state to be preserved")
	}
}
