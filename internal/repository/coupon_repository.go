package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type couponRepository struct {
	db *database.DB
}

func NewCouponRepository(db *database.DB) CouponRepository {
	return &couponRepository{db: db}
}

func (r *couponRepository) FindByCode(code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	if err := r.db.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}
