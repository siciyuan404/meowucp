package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeCouponService struct {
	coupon *domain.Coupon
}

func (f *fakeCouponService) ValidateCoupon(code string, subtotal float64) (*domain.Coupon, error) {
	f.coupon = &domain.Coupon{Code: code, Value: 10}
	return f.coupon, nil
}

func TestCouponValidateEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeCouponService{}
	handler := NewCouponHandler(service)

	r := gin.New()
	r.POST("/api/v1/coupons/validate", handler.Validate)

	body := `{"code":"SAVE10","subtotal":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/coupons/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if service.coupon == nil || service.coupon.Code != "SAVE10" {
		t.Fatalf("expected coupon to be validated")
	}
}
