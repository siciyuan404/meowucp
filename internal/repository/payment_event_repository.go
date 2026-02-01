package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type paymentEventRepository struct {
	db *database.DB
}

func NewPaymentEventRepository(db *database.DB) PaymentEventRepository {
	return &paymentEventRepository{db: db}
}

func (r *paymentEventRepository) Create(event *domain.PaymentEvent) error {
	return r.db.Create(event).Error
}
