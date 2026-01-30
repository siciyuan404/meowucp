package seed

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/meowucp/internal/domain"
)

type fakePaymentHandlerRepo struct {
	items       map[string]*domain.PaymentHandler
	createCount int
}

func newFakePaymentHandlerRepo() *fakePaymentHandlerRepo {
	return &fakePaymentHandlerRepo{items: map[string]*domain.PaymentHandler{}}
}

func (f *fakePaymentHandlerRepo) Create(handler *domain.PaymentHandler) error {
	f.items[handler.Name] = handler
	f.createCount++
	return nil
}

func (f *fakePaymentHandlerRepo) Update(handler *domain.PaymentHandler) error {
	f.items[handler.Name] = handler
	return nil
}

func (f *fakePaymentHandlerRepo) FindByID(id int64) (*domain.PaymentHandler, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePaymentHandlerRepo) FindByName(name string) (*domain.PaymentHandler, error) {
	item, ok := f.items[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (f *fakePaymentHandlerRepo) List() ([]*domain.PaymentHandler, error) {
	result := make([]*domain.PaymentHandler, 0, len(f.items))
	for _, item := range f.items {
		result = append(result, item)
	}
	return result, nil
}

func TestSeedNowPaymentsCreatesWhenMissing(t *testing.T) {
	repo := newFakePaymentHandlerRepo()

	config := NowPaymentsSeedConfig{
		Spec:         "https://nowpayments.io",
		ConfigSchema: "https://nowpayments.io",
		APIBase:      "https://api.nowpayments.io",
		Environment:  "test",
	}

	if err := SeedNowPayments(repo, config); err != nil {
		t.Fatalf("seed nowpayments: %v", err)
	}

	if repo.createCount != 1 {
		t.Fatalf("expected create count 1, got %d", repo.createCount)
	}

	handler, ok := repo.items["com.nowpayments"]
	if !ok {
		t.Fatalf("expected handler to be created")
	}

	if handler.Version == "" {
		t.Fatalf("expected version to be set")
	}

	var parsed map[string]string
	if err := json.Unmarshal([]byte(handler.Config), &parsed); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	if parsed["api_base"] != "https://api.nowpayments.io" {
		t.Fatalf("expected api_base to be set")
	}
	if parsed["environment"] != "test" {
		t.Fatalf("expected environment to be set")
	}
}
