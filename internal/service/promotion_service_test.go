package service

import (
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeCouponRepo struct {
	coupon *domain.Coupon
}

func (f *fakeCouponRepo) FindByCode(code string) (*domain.Coupon, error) {
	return f.coupon, nil
}

func TestCouponValidation(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)
	repo := &fakeCouponRepo{coupon: &domain.Coupon{Code: "SAVE10", Type: "fixed", Value: 10, MinSpend: 50, StartsAt: &start, EndsAt: &end}}
	service := NewPromotionService(repo)

	coupon, err := service.ValidateCoupon("SAVE10", 100)
	if err != nil {
		t.Fatalf("validate coupon: %v", err)
	}
	if coupon == nil || coupon.Code != "SAVE10" {
		t.Fatalf("expected coupon to be returned")
	}
}

func TestPromotionAppliesToTotals(t *testing.T) {
	service := NewPromotionService(nil)
	newTotal, err := service.ApplyPromotions(100, []domain.Promotion{{Name: "ten-off", Rules: "fixed:10"}})
	if err != nil {
		t.Fatalf("apply promotions: %v", err)
	}
	if newTotal != 90 {
		t.Fatalf("expected total 90, got %v", newTotal)
	}
}
