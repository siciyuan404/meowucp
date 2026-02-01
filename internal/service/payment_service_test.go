package service

import (
	"errors"
	"testing"

	"github.com/meowucp/internal/domain"
)

type fakePaymentRepo struct {
	items   []*domain.Payment
	payment *domain.Payment
}

func (f *fakePaymentRepo) Create(payment *domain.Payment) error       { return nil }
func (f *fakePaymentRepo) Update(payment *domain.Payment) error       { f.payment = payment; return nil }
func (f *fakePaymentRepo) FindByID(id int64) (*domain.Payment, error) { return f.payment, nil }
func (f *fakePaymentRepo) FindByOrderID(orderID int64) ([]*domain.Payment, error) {
	if f.payment == nil {
		return []*domain.Payment{}, nil
	}
	return []*domain.Payment{f.payment}, nil
}
func (f *fakePaymentRepo) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	return nil, nil
}
func (f *fakePaymentRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error) {
	if offset >= len(f.items) {
		return []*domain.Payment{}, nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], nil
}
func (f *fakePaymentRepo) Count(filters map[string]interface{}) (int64, error) {
	return int64(len(f.items)), nil
}

type fakePaymentOrderRepo struct {
	order *domain.Order
}

func (f *fakePaymentOrderRepo) Create(order *domain.Order) error                    { return nil }
func (f *fakePaymentOrderRepo) Update(order *domain.Order) error                    { f.order = order; return nil }
func (f *fakePaymentOrderRepo) FindByID(id int64) (*domain.Order, error)            { return f.order, nil }
func (f *fakePaymentOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (f *fakePaymentOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakePaymentOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakePaymentOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakePaymentOrderRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (f *fakePaymentOrderRepo) UpdateStatus(id int64, status string) error {
	if f.order == nil {
		return errors.New("order_not_found")
	}
	f.order.Status = status
	return nil
}
func (f *fakePaymentOrderRepo) CreateOrderItem(item *domain.OrderItem) error { return nil }

type fakePaymentRefundRepo struct {
	created *domain.PaymentRefund
}

func (f *fakePaymentRefundRepo) Create(refund *domain.PaymentRefund) error {
	f.created = refund
	return nil
}

type fakePaymentEventRepo struct {
	created *domain.PaymentEvent
}

func (f *fakePaymentEventRepo) Create(event *domain.PaymentEvent) error {
	f.created = event
	return nil
}

func TestPaymentServiceListPayments(t *testing.T) {
	repo := &fakePaymentRepo{items: []*domain.Payment{{ID: 1}, {ID: 2}, {ID: 3}}}
	svc := NewPaymentService(repo, nil)

	items, total, err := svc.ListPayments(1, 1, map[string]interface{}{})
	if err != nil {
		t.Fatalf("list payments: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(items) != 1 || items[0].ID != 2 {
		t.Fatalf("expected item id 2")
	}
}

func TestPaymentServiceCreateRefundUpdatesPaymentAndOrder(t *testing.T) {
	paymentRepo := &fakePaymentRepo{payment: &domain.Payment{ID: 10, OrderID: 20, Amount: 100, Status: "paid"}}
	orderRepo := &fakePaymentOrderRepo{order: &domain.Order{ID: 20, Status: "paid", PaymentStatus: "paid"}}
	refundRepo := &fakePaymentRefundRepo{}
	eventRepo := &fakePaymentEventRepo{}

	service := NewPaymentServiceWithDeps(paymentRepo, orderRepo, refundRepo, eventRepo)
	refund, err := service.CreateRefund(10, 40, "customer_request")
	if err != nil {
		t.Fatalf("create refund: %v", err)
	}
	if refund == nil || refund.PaymentID != 10 {
		t.Fatalf("expected refund record for payment")
	}
	if paymentRepo.payment.Status != "partially_refunded" {
		t.Fatalf("expected payment status partially_refunded")
	}
	if orderRepo.order.Status != "paid" {
		t.Fatalf("expected order to remain paid")
	}
	if refundRepo.created == nil {
		t.Fatalf("expected refund repo to be called")
	}
	if eventRepo.created == nil {
		t.Fatalf("expected payment event to be recorded")
	}
}
