package service

import (
	"errors"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type UCPOrderService struct {
	orderRepo   repository.OrderRepository
	paymentRepo repository.PaymentRepository
}

func NewUCPOrderService(orderRepo repository.OrderRepository, paymentRepo repository.PaymentRepository) *UCPOrderService {
	return &UCPOrderService{
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
	}
}

type ucpOrderTransactionRunner interface {
	Transaction(fn func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, paymentRepo repository.PaymentRepository) error) error
}

func (s *UCPOrderService) CreateFromCheckout(order *domain.Order, items []*domain.OrderItem, payment *domain.Payment) (*domain.Order, *domain.Payment, error) {
	if txRunner, ok := s.orderRepo.(ucpOrderTransactionRunner); ok {
		var createdOrder *domain.Order
		var createdPayment *domain.Payment
		err := txRunner.Transaction(func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, paymentRepo repository.PaymentRepository) error {
			var err error
			createdOrder, createdPayment, err = s.createFromCheckoutWithRepos(orderRepo, paymentRepo, order, items, payment)
			return err
		})
		if err != nil {
			return nil, nil, err
		}
		return createdOrder, createdPayment, nil
	}

	return s.createFromCheckoutWithRepos(s.orderRepo, s.paymentRepo, order, items, payment)
}

func (s *UCPOrderService) createFromCheckoutWithRepos(orderRepo repository.OrderRepository, paymentRepo repository.PaymentRepository, order *domain.Order, items []*domain.OrderItem, payment *domain.Payment) (*domain.Order, *domain.Payment, error) {
	if orderRepo == nil {
		return nil, nil, errors.New("order repository unavailable")
	}
	if payment != nil && paymentRepo == nil {
		return nil, nil, errors.New("payment repository unavailable")
	}
	if err := orderRepo.Create(order); err != nil {
		return nil, nil, err
	}

	for _, item := range items {
		item.OrderID = order.ID
		if err := orderRepo.CreateOrderItem(item); err != nil {
			return nil, nil, err
		}
	}

	if payment != nil {
		payment.OrderID = order.ID
		if err := paymentRepo.Create(payment); err != nil {
			return nil, nil, err
		}
	}

	return order, payment, nil
}
