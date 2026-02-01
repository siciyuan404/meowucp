package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type PromotionService struct {
	couponRepo repository.CouponRepository
}

func NewPromotionService(repo repository.CouponRepository) *PromotionService {
	return &PromotionService{couponRepo: repo}
}

func (s *PromotionService) ValidateCoupon(code string, subtotal float64) (*domain.Coupon, error) {
	if s == nil || s.couponRepo == nil {
		return nil, errors.New("coupon_repository_unavailable")
	}
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("coupon_code_required")
	}
	coupon, err := s.couponRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if coupon.StartsAt != nil && now.Before(*coupon.StartsAt) {
		return nil, errors.New("coupon_not_started")
	}
	if coupon.EndsAt != nil && now.After(*coupon.EndsAt) {
		return nil, errors.New("coupon_expired")
	}
	if subtotal < coupon.MinSpend {
		return nil, errors.New("coupon_min_spend")
	}
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return nil, errors.New("coupon_usage_limit")
	}
	return coupon, nil
}

func (s *PromotionService) ApplyPromotions(subtotal float64, promotions []domain.Promotion) (float64, error) {
	newTotal := subtotal
	for _, promo := range promotions {
		value := strings.TrimSpace(promo.Rules)
		if strings.HasPrefix(value, "fixed:") {
			off, err := strconv.ParseFloat(strings.TrimPrefix(value, "fixed:"), 64)
			if err != nil {
				return subtotal, err
			}
			newTotal -= off
		}
	}
	if newTotal < 0 {
		newTotal = 0
	}
	return newTotal, nil
}
