package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (s *PaymentService) CreatePayment(payment *domain.Payment) error {
	return s.paymentRepo.Create(payment)
}

func (s *PaymentService) ProcessPayment(orderID int64, amount float64, paymentMethod string) (*domain.Payment, error) {
	_, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	payment := &domain.Payment{
		OrderID:       orderID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "pending",
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPayment(id int64) (*domain.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *PaymentService) GetOrderPayments(orderID int64) ([]*domain.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}

func (s *PaymentService) ListPayments(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, int64, error) {
	items, err := s.paymentRepo.List(offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.paymentRepo.Count(filters)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
