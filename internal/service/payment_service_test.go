package service

import (
	"testing"

	"github.com/meowucp/internal/domain"
)

type fakePaymentRepo struct {
	items []*domain.Payment
}

func (f *fakePaymentRepo) Create(payment *domain.Payment) error                   { return nil }
func (f *fakePaymentRepo) Update(payment *domain.Payment) error                   { return nil }
func (f *fakePaymentRepo) FindByID(id int64) (*domain.Payment, error)             { return nil, nil }
func (f *fakePaymentRepo) FindByOrderID(orderID int64) ([]*domain.Payment, error) { return nil, nil }
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
