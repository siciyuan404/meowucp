package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeShippingRateService struct {
	tax      float64
	shipping float64
}

func (f *fakeShippingRateService) Quote(region string, items []domain.OrderItem) (float64, float64, error) {
	return f.tax, f.shipping, nil
}

func TestShippingRateEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeShippingRateService{tax: 1.5, shipping: 10}
	handler := NewShippingRateHandler(service)

	r := gin.New()
	r.GET("/api/v1/shipping/rates", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/shipping/rates?region=CN&quantity=2", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}

func TestAddressValidationEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAddressValidationHandler()

	r := gin.New()
	r.POST("/api/v1/address/validate", handler.Validate)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/address/validate", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
