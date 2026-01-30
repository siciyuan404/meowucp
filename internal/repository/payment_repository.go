package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type paymentRepository struct {
	db *database.DB
}

func NewPaymentRepository(db *database.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) FindByID(id int64) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderID(orderID int64) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.Where("order_id = ?", orderID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	query := r.db.Model(&domain.Payment{})
	for key, value := range filters {
		query = query.Where(key, value)
	}
	err := query.Offset(offset).Limit(limit).Order("id desc").Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) Count(filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&domain.Payment{})
	for key, value := range filters {
		query = query.Where(key, value)
	}
	err := query.Count(&count).Error
	return count, err
}
