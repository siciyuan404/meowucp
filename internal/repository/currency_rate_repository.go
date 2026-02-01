package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type currencyRateRepository struct {
	db *database.DB
}

func NewCurrencyRateRepository(db *database.DB) CurrencyRateRepository {
	return &currencyRateRepository{db: db}
}

func (r *currencyRateRepository) FindByBaseAndTarget(base, target string) (*domain.CurrencyRate, error) {
	var rate domain.CurrencyRate
	if err := r.db.Where("base = ? AND target = ?", base, target).First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *currencyRateRepository) Create(rate *domain.CurrencyRate) error {
	return r.db.Create(rate).Error
}
