package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type shippingRuleRepository struct {
	db *database.DB
}

func NewShippingRuleRepository(db *database.DB) ShippingRuleRepository {
	return &shippingRuleRepository{db: db}
}

func (r *shippingRuleRepository) ListByRegion(region string) ([]*domain.ShippingRule, error) {
	items := []*domain.ShippingRule{}
	if err := r.db.Where("region = ?", region).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
