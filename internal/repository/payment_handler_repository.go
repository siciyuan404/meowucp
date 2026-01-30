package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type paymentHandlerRepository struct {
	db *database.DB
}

func NewPaymentHandlerRepository(db *database.DB) PaymentHandlerRepository {
	return &paymentHandlerRepository{db: db}
}

func (r *paymentHandlerRepository) Create(handler *domain.PaymentHandler) error {
	return r.db.Create(handler).Error
}

func (r *paymentHandlerRepository) Update(handler *domain.PaymentHandler) error {
	return r.db.Save(handler).Error
}

func (r *paymentHandlerRepository) FindByID(id int64) (*domain.PaymentHandler, error) {
	var handler domain.PaymentHandler
	err := r.db.First(&handler, id).Error
	if err != nil {
		return nil, err
	}
	return &handler, nil
}

func (r *paymentHandlerRepository) FindByName(name string) (*domain.PaymentHandler, error) {
	var handler domain.PaymentHandler
	err := r.db.Where("name = ?", name).First(&handler).Error
	if err != nil {
		return nil, err
	}
	return &handler, nil
}

func (r *paymentHandlerRepository) List() ([]*domain.PaymentHandler, error) {
	var handlers []*domain.PaymentHandler
	err := r.db.Find(&handlers).Error
	return handlers, err
}
