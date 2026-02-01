package service

import (
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
)

type fakeTaxRuleRepo struct {
	rules []*domain.TaxRule
}

func (f *fakeTaxRuleRepo) ListByRegion(region string) ([]*domain.TaxRule, error) {
	return f.rules, nil
}

type fakeShippingRuleRepo struct {
	rules []*domain.ShippingRule
}

func (f *fakeShippingRuleRepo) ListByRegion(region string) ([]*domain.ShippingRule, error) {
	return f.rules, nil
}

func TestTaxRateAppliedByRegion(t *testing.T) {
	taxRepo := &fakeTaxRuleRepo{rules: []*domain.TaxRule{{Region: "CN", Rate: 0.1, EffectiveAt: time.Now().Add(-time.Hour)}}}
	shippingRepo := &fakeShippingRuleRepo{}
	service := NewTaxShippingService(taxRepo, shippingRepo)

	items := []domain.OrderItem{{Quantity: 2, UnitPrice: 50, TotalPrice: 100}}
	tax, _, err := service.Quote("CN", items)
	if err != nil {
		t.Fatalf("quote: %v", err)
	}
	if tax != 10 {
		t.Fatalf("expected tax 10, got %v", tax)
	}
}

func TestShippingRateByItems(t *testing.T) {
	taxRepo := &fakeTaxRuleRepo{}
	shippingRepo := &fakeShippingRuleRepo{rules: []*domain.ShippingRule{{Region: "CN", BaseAmount: 5, PerItemAmount: 2}}}
	service := NewTaxShippingService(taxRepo, shippingRepo)

	items := []domain.OrderItem{{Quantity: 3, UnitPrice: 10, TotalPrice: 30}}
	_, shipping, err := service.Quote("CN", items)
	if err != nil {
		t.Fatalf("quote: %v", err)
	}
	if shipping != 11 {
		t.Fatalf("expected shipping 11, got %v", shipping)
	}
}
