package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type paymentRefundRepository struct {
	db *database.DB
}

func NewPaymentRefundRepository(db *database.DB) PaymentRefundRepository {
	return &paymentRefundRepository{db: db}
}

func (r *paymentRefundRepository) Create(refund *domain.PaymentRefund) error {
	return r.db.Create(refund).Error
}
