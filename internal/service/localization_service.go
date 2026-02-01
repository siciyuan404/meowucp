package service

import (
	"github.com/meowucp/internal/repository"
)

type LocalizationService struct {
	currencyRateRepo repository.CurrencyRateRepository
	i18nStringRepo   repository.I18nStringRepository
}

func NewLocalizationService(currencyRateRepo repository.CurrencyRateRepository, i18nStringRepo repository.I18nStringRepository) *LocalizationService {
	return &LocalizationService{
		currencyRateRepo: currencyRateRepo,
		i18nStringRepo:   i18nStringRepo,
	}
}

func (s *LocalizationService) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}
	rate, err := s.currencyRateRepo.FindByBaseAndTarget(base, target)
	if err != nil {
		return 0, err
	}
	return amount * rate.Rate, nil
}

func (s *LocalizationService) Translate(key, locale string) (string, error) {
	str, err := s.i18nStringRepo.FindByKeyAndLocale(key, locale)
	if err != nil {
		return "", err
	}
	if str == nil {
		return key, nil
	}
	return str.Value, nil
}
