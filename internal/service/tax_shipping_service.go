package service

import (
	"errors"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type TaxShippingService struct {
	taxRepo      repository.TaxRuleRepository
	shippingRepo repository.ShippingRuleRepository
}

func NewTaxShippingService(taxRepo repository.TaxRuleRepository, shippingRepo repository.ShippingRuleRepository) *TaxShippingService {
	return &TaxShippingService{taxRepo: taxRepo, shippingRepo: shippingRepo}
}

func (s *TaxShippingService) Quote(region string, items []domain.OrderItem) (float64, float64, error) {
	if s == nil {
		return 0, 0, errors.New("service_unavailable")
	}
	if region == "" {
		return 0, 0, errors.New("region_required")
	}
	subtotal := 0.0
	quantity := 0
	for _, item := range items {
		if item.TotalPrice > 0 {
			subtotal += item.TotalPrice
		} else {
			subtotal += item.UnitPrice * float64(item.Quantity)
		}
		quantity += item.Quantity
	}

	taxRate := 0.0
	if s.taxRepo != nil {
		rules, err := s.taxRepo.ListByRegion(region)
		if err != nil {
			return 0, 0, err
		}
		if len(rules) > 0 {
			selected := selectLatestTaxRule(rules)
			if selected != nil {
				taxRate = selected.Rate
			}
		}
	}

	shipping := 0.0
	if s.shippingRepo != nil {
		rules, err := s.shippingRepo.ListByRegion(region)
		if err != nil {
			return 0, 0, err
		}
		if len(rules) > 0 {
			rule := rules[0]
			shipping = rule.BaseAmount + rule.PerItemAmount*float64(quantity)
		}
	}

	return subtotal * taxRate, shipping, nil
}

func selectLatestTaxRule(rules []*domain.TaxRule) *domain.TaxRule {
	if len(rules) == 0 {
		return nil
	}
	latest := rules[0]
	for _, rule := range rules[1:] {
		if rule == nil {
			continue
		}
		if rule.EffectiveAt.After(latest.EffectiveAt) || latest.EffectiveAt.IsZero() {
			latest = rule
		}
	}
	if latest == nil {
		return nil
	}
	if latest.EffectiveAt.After(time.Now()) {
		return nil
	}
	return latest
}
