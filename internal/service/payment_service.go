package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	refundRepo  repository.PaymentRefundRepository
	eventRepo   repository.PaymentEventRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func NewPaymentServiceWithDeps(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, refundRepo repository.PaymentRefundRepository, eventRepo repository.PaymentEventRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		refundRepo:  refundRepo,
		eventRepo:   eventRepo,
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

func (s *PaymentService) MarkPaymentPaid(orderID int64, transactionID string) error {
	if s == nil || s.paymentRepo == nil {
		return errors.New("payment repository unavailable")
	}
	payments, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return err
	}
	if len(payments) == 0 {
		return errors.New("payment not found")
	}
	payment := payments[0]
	payment.Status = "paid"
	if transactionID != "" {
		payment.TransactionID = transactionID
	}
	return s.paymentRepo.Update(payment)
}

func (s *PaymentService) CreateRefund(paymentID int64, amount float64, reason string) (*domain.PaymentRefund, error) {
	if s == nil || s.paymentRepo == nil || s.orderRepo == nil || s.refundRepo == nil || s.eventRepo == nil {
		return nil, errors.New("refund_dependencies_unavailable")
	}
	if amount <= 0 {
		return nil, errors.New("invalid_refund_amount")
	}
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil || payment == nil {
		return nil, errors.New("payment_not_found")
	}
	if amount > payment.Amount {
		return nil, errors.New("refund_amount_exceeds_payment")
	}

	refund := &domain.PaymentRefund{
		PaymentID: paymentID,
		Amount:    amount,
		Status:    "completed",
		Reason:    reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.refundRepo.Create(refund); err != nil {
		return nil, err
	}

	status := "partially_refunded"
	if amount == payment.Amount {
		status = "refunded"
	}
	payment.Status = status
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	if status == "refunded" {
		order, err := s.orderRepo.FindByID(payment.OrderID)
		if err != nil || order == nil {
			return nil, errors.New("order_not_found")
		}
		now := time.Now()
		order.Status = "refunded"
		order.PaymentStatus = "refunded"
		order.RefundedAt = &now
		if err := s.orderRepo.Update(order); err != nil {
			return nil, err
		}
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"payment_id": paymentID,
		"amount":     amount,
		"reason":     reason,
	})
	payloadText := string(payload)
	if err := s.eventRepo.Create(&domain.PaymentEvent{
		PaymentID: paymentID,
		EventType: "refund_created",
		Payload:   &payloadText,
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return refund, nil
}
