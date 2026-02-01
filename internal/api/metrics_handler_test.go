package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMetricsEndpointRequiresInternalAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewMetricsHandler("metrics-token")
	r := gin.New()
	r.GET("/metrics", handler.Serve)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}
