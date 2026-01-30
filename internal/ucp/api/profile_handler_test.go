package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

type fakePaymentHandlerRepo struct {
	items map[string]*domain.PaymentHandler
}

func newFakePaymentHandlerRepo() *fakePaymentHandlerRepo {
	return &fakePaymentHandlerRepo{items: map[string]*domain.PaymentHandler{}}
}

func (f *fakePaymentHandlerRepo) Create(handler *domain.PaymentHandler) error {
	f.items[handler.Name] = handler
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

var _ repository.PaymentHandlerRepository = (*fakePaymentHandlerRepo)(nil)

func TestProfileIncludesNowPaymentsInstrumentSchema(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakePaymentHandlerRepo()
	configJSON, _ := json.Marshal(map[string]string{
		"api_base": "https://api.nowpayments.io",
	})
	repo.Create(&domain.PaymentHandler{
		Name:         "com.nowpayments",
		Version:      "2026-01-11",
		Spec:         "https://nowpayments.io",
		ConfigSchema: "https://nowpayments.io",
		Config:       string(configJSON),
	})

	services := &service.Services{Handler: service.NewPaymentHandlerService(repo)}
	handler := NewProfileHandler(services)

	r := gin.New()
	r.GET("/.well-known/ucp", handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/ucp", nil)
	req.Host = "example.com"
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var profile model.Profile
	if err := json.Unmarshal(resp.Body.Bytes(), &profile); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if profile.Payment == nil || len(profile.Payment.Handlers) != 1 {
		t.Fatalf("expected one payment handler")
	}

	handlerResp := profile.Payment.Handlers[0]
	if len(handlerResp.InstrumentSchemas) == 0 {
		t.Fatalf("expected instrument_schemas to be set")
	}

	configMap, ok := handlerResp.Config.(map[string]interface{})
	if !ok {
		t.Fatalf("expected config to be decoded json object")
	}
	if configMap["api_base"] != "https://api.nowpayments.io" {
		t.Fatalf("expected api_base in config")
	}
}
