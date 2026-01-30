package service

import (
	"errors"
	"testing"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type ucpOrderStore struct {
	orders   []*domain.Order
	items    []*domain.OrderItem
	payments []*domain.Payment
}

func (s *ucpOrderStore) clone() *ucpOrderStore {
	clone := &ucpOrderStore{
		orders:   append([]*domain.Order{}, s.orders...),
		items:    append([]*domain.OrderItem{}, s.items...),
		payments: append([]*domain.Payment{}, s.payments...),
	}
	return clone
}

type ucpFakeOrderRepo struct {
	store *ucpOrderStore
}

func (r *ucpFakeOrderRepo) Create(order *domain.Order) error {
	order.ID = 1
	r.store.orders = append(r.store.orders, order)
	return nil
}
func (r *ucpFakeOrderRepo) Update(order *domain.Order) error                    { return nil }
func (r *ucpFakeOrderRepo) FindByID(id int64) (*domain.Order, error)            { return nil, nil }
func (r *ucpFakeOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (r *ucpFakeOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (r *ucpFakeOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (r *ucpFakeOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (r *ucpFakeOrderRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *ucpFakeOrderRepo) UpdateStatus(id int64, status string) error          { return nil }
func (r *ucpFakeOrderRepo) CreateOrderItem(item *domain.OrderItem) error {
	r.store.items = append(r.store.items, item)
	return nil
}

type ucpTxOrderRepo struct {
	store             *ucpOrderStore
	transactionCalled bool
	paymentCreateErr  error
}

func (r *ucpTxOrderRepo) Create(order *domain.Order) error {
	order.ID = 1
	r.store.orders = append(r.store.orders, order)
	return nil
}
func (r *ucpTxOrderRepo) Update(order *domain.Order) error                    { return nil }
func (r *ucpTxOrderRepo) FindByID(id int64) (*domain.Order, error)            { return nil, nil }
func (r *ucpTxOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (r *ucpTxOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (r *ucpTxOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (r *ucpTxOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (r *ucpTxOrderRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *ucpTxOrderRepo) UpdateStatus(id int64, status string) error          { return nil }
func (r *ucpTxOrderRepo) CreateOrderItem(item *domain.OrderItem) error {
	r.store.items = append(r.store.items, item)
	return nil
}
func (r *ucpTxOrderRepo) Transaction(fn func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, paymentRepo repository.PaymentRepository) error) error {
	r.transactionCalled = true
	clone := r.store.clone()
	orderRepo := &ucpTxOrderRepo{store: clone}
	paymentRepo := &ucpTxPaymentRepo{store: clone, createErr: r.paymentCreateErr}
	if err := fn(orderRepo, nil, nil, nil, paymentRepo); err != nil {
		return err
	}
	*r.store = *clone
	return nil
}

type ucpFakePaymentRepo struct {
	store     *ucpOrderStore
	createErr error
}

func (r *ucpFakePaymentRepo) Create(payment *domain.Payment) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.store.payments = append(r.store.payments, payment)
	return nil
}
func (r *ucpFakePaymentRepo) Update(payment *domain.Payment) error { return nil }
func (r *ucpFakePaymentRepo) FindByID(id int64) (*domain.Payment, error) {
	return nil, nil
}
func (r *ucpFakePaymentRepo) FindByOrderID(orderID int64) ([]*domain.Payment, error) {
	return nil, nil
}
func (r *ucpFakePaymentRepo) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	return nil, nil
}
func (r *ucpFakePaymentRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error) {
	return nil, nil
}
func (r *ucpFakePaymentRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }

type ucpTxPaymentRepo struct {
	store     *ucpOrderStore
	createErr error
}

func (r *ucpTxPaymentRepo) Create(payment *domain.Payment) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.store.payments = append(r.store.payments, payment)
	return nil
}
func (r *ucpTxPaymentRepo) Update(payment *domain.Payment) error { return nil }
func (r *ucpTxPaymentRepo) FindByID(id int64) (*domain.Payment, error) {
	return nil, nil
}
func (r *ucpTxPaymentRepo) FindByOrderID(orderID int64) ([]*domain.Payment, error) {
	return nil, nil
}
func (r *ucpTxPaymentRepo) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	return nil, nil
}
func (r *ucpTxPaymentRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error) {
	return nil, nil
}
func (r *ucpTxPaymentRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }

func TestUCPOrderServiceCreateFromCheckoutRollsBackOnPaymentFailure(t *testing.T) {
	store := &ucpOrderStore{}
	orderRepo := &ucpTxOrderRepo{store: store, paymentCreateErr: errors.New("payment failed")}
	paymentRepo := &ucpTxPaymentRepo{store: store, createErr: errors.New("payment failed")}
	svc := NewUCPOrderService(orderRepo, paymentRepo)

	order := &domain.Order{}
	items := []*domain.OrderItem{{
		ProductName: "Cat Toy",
		SKU:         "CAT-TOY-001",
		Quantity:    1,
		UnitPrice:   10,
		TotalPrice:  10,
	}}
	payment := &domain.Payment{Amount: 10, PaymentMethod: "card"}

	_, _, err := svc.CreateFromCheckout(order, items, payment)
	if err == nil {
		t.Fatalf("expected error when payment create fails")
	}
	if !orderRepo.transactionCalled {
		t.Fatalf("expected transaction to be used")
	}
	if len(store.orders) != 0 || len(store.items) != 0 || len(store.payments) != 0 {
		t.Fatalf("expected no persisted order, items, or payments after failure")
	}
}
