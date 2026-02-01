package service

import (
	"errors"
	"testing"

	"github.com/meowucp/internal/domain"
)

type fakeCurrencyRateRepository struct {
	rates map[string]*domain.CurrencyRate
}

func newFakeCurrencyRateRepository() *fakeCurrencyRateRepository {
	return &fakeCurrencyRateRepository{
		rates: map[string]*domain.CurrencyRate{
			"USD-CNY": {ID: 1, Base: "USD", Target: "CNY", Rate: 7.2},
			"USD-EUR": {ID: 2, Base: "USD", Target: "EUR", Rate: 0.92},
		},
	}
}

func (f *fakeCurrencyRateRepository) FindByBaseAndTarget(base, target string) (*domain.CurrencyRate, error) {
	key := base + "-" + target
	if rate, ok := f.rates[key]; ok {
		return rate, nil
	}
	return nil, errors.New("rate not found")
}

func (f *fakeCurrencyRateRepository) Create(rate *domain.CurrencyRate) error {
	f.rates[rate.Base+"-"+rate.Target] = rate
	return nil
}

type fakeI18nStringRepository struct {
	strings map[string]*domain.I18nString
}

func newFakeI18nStringRepository() *fakeI18nStringRepository {
	return &fakeI18nStringRepository{
		strings: map[string]*domain.I18nString{
			"product.name.en-US": {ID: 1, Key: "product.name", Locale: "en-US", Value: "Product"},
			"product.name.zh-CN": {ID: 2, Key: "product.name", Locale: "zh-CN", Value: "产品"},
		},
	}
}

func (f *fakeI18nStringRepository) FindByKeyAndLocale(key, locale string) (*domain.I18nString, error) {
	if str, ok := f.strings[key+"."+locale]; ok {
		return str, nil
	}
	return nil, nil
}

func (f *fakeI18nStringRepository) Create(str *domain.I18nString) error {
	f.strings[str.Key+"."+str.Locale] = str
	return nil
}

func TestConvertCurrency(t *testing.T) {
	svc := NewLocalizationService(newFakeCurrencyRateRepository(), newFakeI18nStringRepository())

	tests := []struct {
		input    float64
		base     string
		target   string
		expected float64
	}{
		{100, "USD", "USD", 100},
		{100, "USD", "CNY", 720},
		{100, "USD", "EUR", 92},
	}

	for _, tt := range tests {
		result, err := svc.Convert(tt.input, tt.base, tt.target)
		if err != nil {
			t.Errorf("Convert(%f, %s, %s) error: %v", tt.input, tt.base, tt.target, err)
		}
		if result != tt.expected {
			t.Errorf("Convert(%f, %s, %s) = %f, expected %f", tt.input, tt.base, tt.target, result, tt.expected)
		}
	}
}

func TestTranslateString(t *testing.T) {
	svc := NewLocalizationService(newFakeCurrencyRateRepository(), newFakeI18nStringRepository())

	tests := []struct {
		key      string
		locale   string
		expected string
	}{
		{"product.name", "en-US", "Product"},
		{"product.name", "zh-CN", "产品"},
		{"missing.key", "en-US", "missing.key"},
	}

	for _, tt := range tests {
		result, err := svc.Translate(tt.key, tt.locale)
		if err != nil {
			t.Errorf("Translate(%s, %s) error: %v", tt.key, tt.locale, err)
		}
		if result != tt.expected {
			t.Errorf("Translate(%s, %s) = %s, expected %s", tt.key, tt.locale, result, tt.expected)
		}
	}
}
