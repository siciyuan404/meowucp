package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type taxRuleRepository struct {
	db *database.DB
}

func NewTaxRuleRepository(db *database.DB) TaxRuleRepository {
	return &taxRuleRepository{db: db}
}

func (r *taxRuleRepository) ListByRegion(region string) ([]*domain.TaxRule, error) {
	items := []*domain.TaxRule{}
	if err := r.db.Where("region = ?", region).Order("effective_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
