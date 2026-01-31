package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

type fakeCheckoutRepo struct {
	items map[string]*domain.CheckoutSession
}

func newFakeCheckoutRepo() *fakeCheckoutRepo {
	return &fakeCheckoutRepo{items: map[string]*domain.CheckoutSession{}}
}

func (f *fakeCheckoutRepo) Create(session *domain.CheckoutSession) error {
	f.items[session.ID] = session
	return nil
}

func (f *fakeCheckoutRepo) Update(session *domain.CheckoutSession) error {
	f.items[session.ID] = session
	return nil
}

func (f *fakeCheckoutRepo) FindByID(id string) (*domain.CheckoutSession, error) {
	session, ok := f.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return session, nil
}

func (f *fakeCheckoutRepo) Delete(id string) error {
	delete(f.items, id)
	return nil
}

type fakeOrderRepo struct {
	orders      map[int64]*domain.Order
	orderItems  []*domain.OrderItem
	createCount int
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{orders: map[int64]*domain.Order{}}
}

func (f *fakeOrderRepo) Create(order *domain.Order) error {
	order.ID = int64(len(f.orders) + 1)
	f.orders[order.ID] = order
	f.createCount++
	return nil
}

func (f *fakeOrderRepo) Update(order *domain.Order) error {
	f.orders[order.ID] = order
	return nil
}

func (f *fakeOrderRepo) FindByID(id int64) (*domain.Order, error) {
	order, ok := f.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return order, nil
}

func (f *fakeOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepo) CountByUserID(userID int64) (int64, error) {
	return 0, errors.New("not implemented")
}

func (f *fakeOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepo) Count(filters map[string]interface{}) (int64, error) {
	return 0, errors.New("not implemented")
}

func (f *fakeOrderRepo) UpdateStatus(id int64, status string) error {
	order, ok := f.orders[id]
	if !ok {
		return errors.New("not found")
	}
	order.Status = status
	return nil
}

func (f *fakeOrderRepo) CreateOrderItem(item *domain.OrderItem) error {
	f.orderItems = append(f.orderItems, item)
	return nil
}

type fakePaymentRepo struct {
	items       []*domain.Payment
	createCount int
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{}
}

func (f *fakePaymentRepo) Create(payment *domain.Payment) error {
	payment.ID = int64(len(f.items) + 1)
	f.items = append(f.items, payment)
	f.createCount++
	return nil
}

func (f *fakePaymentRepo) Update(payment *domain.Payment) error {
	return nil
}

func (f *fakePaymentRepo) FindByID(id int64) (*domain.Payment, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePaymentRepo) FindByOrderID(orderID int64) ([]*domain.Payment, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePaymentRepo) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	return nil, errors.New("not implemented")
}

func (f *fakePaymentRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error) {
	return f.items, nil
}

func (f *fakePaymentRepo) Count(filters map[string]interface{}) (int64, error) {
	return int64(len(f.items)), nil
}

func TestCheckoutCreateAndGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(repo)
	paymentHandlerRepo := newFakePaymentHandlerRepo()
	configJSON, _ := json.Marshal(map[string]string{
		"api_base": "https://api.nowpayments.io",
	})
	paymentHandlerRepo.Create(&domain.PaymentHandler{
		Name:         "com.nowpayments",
		Version:      "2026-01-11",
		Spec:         "https://nowpayments.io",
		ConfigSchema: "https://nowpayments.io",
		Config:       string(configJSON),
	})
	services := &service.Services{
		Checkout: checkoutService,
		Handler:  service.NewPaymentHandlerService(paymentHandlerRepo),
	}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)
	r.GET("/ucp/v1/checkout-sessions/:id", handler.Get)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if created.Currency != "CNY" {
		t.Fatalf("expected currency CNY, got %s", created.Currency)
	}
	if created.Status != "ready_for_complete" {
		t.Fatalf("expected status ready_for_complete, got %s", created.Status)
	}
	if len(created.Totals) == 0 {
		t.Fatalf("expected totals to be set")
	}
	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/ucp/v1/checkout-sessions/"+created.ID, nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getResp.Code)
	}

	var fetched model.CheckoutSession
	if err := json.Unmarshal(getResp.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("unmarshal get response: %v", err)
	}
	if fetched.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, fetched.ID)
	}
}

func TestCheckoutCreateIncludesPaymentHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	paymentHandlerRepo := newFakePaymentHandlerRepo()
	configJSON, _ := json.Marshal(map[string]string{
		"api_base": "https://api.nowpayments.io",
	})
	paymentHandlerRepo.Create(&domain.PaymentHandler{
		Name:         "com.nowpayments",
		Version:      "2026-01-11",
		Spec:         "https://nowpayments.io",
		ConfigSchema: "https://nowpayments.io",
		Config:       string(configJSON),
	})

	services := &service.Services{
		Checkout: checkoutService,
		Handler:  service.NewPaymentHandlerService(paymentHandlerRepo),
	}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(created.Payment.Handlers) != 1 {
		t.Fatalf("expected one payment handler")
	}
	if len(created.Payment.Handlers[0].InstrumentSchemas) == 0 {
		t.Fatalf("expected instrument_schemas to be set")
	}
	if _, ok := created.Payment.Handlers[0].Config.(map[string]interface{}); !ok {
		t.Fatalf("expected config to be decoded json object")
	}
}

func TestCheckoutCreateIncludesLinksAndContinueURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	services := &service.Services{Checkout: checkoutService}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Host = "example.com"
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(created.Links) == 0 {
		t.Fatalf("expected links to be set")
	}
	if created.ContinueURL == "" {
		t.Fatalf("expected continue_url to be set")
	}
}

func TestCheckoutCreateUsesConfiguredLinksAndContinueURLBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	services := &service.Services{Checkout: checkoutService}

	config := CheckoutHandlerConfig{
		Links: []model.Link{
			{Type: "privacy_policy", URL: "https://merchant.example.com/privacy"},
			{Type: "terms_of_service", URL: "https://merchant.example.com/terms"},
		},
		ContinueURLBase: "https://pay.example.com/checkout-sessions",
	}

	handler := NewCheckoutHandlerWithConfig(services, config)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(created.Links) != 2 {
		t.Fatalf("expected configured links to be used")
	}
	if created.ContinueURL == "" || !strings.HasPrefix(created.ContinueURL, config.ContinueURLBase) {
		t.Fatalf("expected continue_url to use configured base")
	}
}

func TestCheckoutCreateRequiresEscalationWhenNoHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	services := &service.Services{Checkout: checkoutService}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if created.Status != "requires_escalation" {
		t.Fatalf("expected status requires_escalation, got %s", created.Status)
	}
	if created.ContinueURL == "" {
		t.Fatalf("expected continue_url to be set")
	}
	if len(created.Messages) == 0 {
		t.Fatalf("expected messages to be set")
	}

	found := false
	for _, message := range created.Messages {
		if message.Code == "payment_required" && message.Severity == "requires_buyer_input" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected payment_required requires_buyer_input message")
	}
}

func TestCheckoutCreateRequiresEscalationOnRecoverableErrorsWithoutHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	services := &service.Services{Checkout: checkoutService}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)

	tests := map[string]model.CheckoutCreateRequest{
		"missing_currency": {
			LineItems: []model.LineItem{
				{
					Item: model.Item{
						ID:    "sku_1",
						Title: "Test Item",
						Price: 19900,
					},
					Quantity: 1,
				},
			},
		},
		"missing_line_items": {
			Currency: "CNY",
		},
	}

	for name, createBody := range tests {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(createBody)
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}

			createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
			createReq.Header.Set("Content-Type", "application/json")
			createResp := httptest.NewRecorder()
			r.ServeHTTP(createResp, createReq)

			if createResp.Code != http.StatusCreated {
				t.Fatalf("expected status 201, got %d", createResp.Code)
			}

			var created model.CheckoutSession
			if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
				t.Fatalf("unmarshal response: %v", err)
			}

			if created.Status != "requires_escalation" {
				t.Fatalf("expected status requires_escalation, got %s", created.Status)
			}

			recoverable := false
			for _, message := range created.Messages {
				if message.Severity == "recoverable" {
					recoverable = true
					break
				}
			}
			if !recoverable {
				t.Fatalf("expected recoverable message")
			}
		})
	}
}

func TestCheckoutUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(repo)
	services := &service.Services{Checkout: checkoutService}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)
	r.PUT("/ucp/v1/checkout-sessions/:id", handler.Update)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal create request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	updateBody := model.CheckoutUpdateRequest{
		ID:       created.ID,
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 2,
			},
		},
	}

	updatePayload, err := json.Marshal(updateBody)
	if err != nil {
		t.Fatalf("marshal update request: %v", err)
	}

	updateReq := httptest.NewRequest(http.MethodPut, "/ucp/v1/checkout-sessions/"+created.ID, bytes.NewReader(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	r.ServeHTTP(updateResp, updateReq)

	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", updateResp.Code)
	}

	var updated model.CheckoutSession
	if err := json.Unmarshal(updateResp.Body.Bytes(), &updated); err != nil {
		t.Fatalf("unmarshal update response: %v", err)
	}

	if len(updated.LineItems) != 1 || updated.LineItems[0].Quantity != 2 {
		t.Fatalf("expected quantity updated to 2")
	}
	if len(updated.Totals) == 0 || updated.Totals[0].Amount != 39800 {
		t.Fatalf("expected totals to reflect updated quantity")
	}
}

func TestCheckoutCompleteCreatesOrderAndPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	orderRepo := newFakeOrderRepo()
	paymentRepo := newFakePaymentRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	ucpOrderService := service.NewUCPOrderService(orderRepo, paymentRepo)

	services := &service.Services{
		Checkout: checkoutService,
		UCPOrder: ucpOrderService,
	}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)
	r.POST("/ucp/v1/checkout-sessions/:id/complete", handler.Complete)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal create request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	completeBody := model.CheckoutCompleteRequest{
		PaymentData: model.PaymentInstrument{
			HandlerID: "com.nowpayments",
			Type:      "card",
			Credential: map[string]string{
				"token": "tok_123",
			},
		},
	}

	completePayload, err := json.Marshal(completeBody)
	if err != nil {
		t.Fatalf("marshal complete request: %v", err)
	}

	completeReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions/"+created.ID+"/complete", bytes.NewReader(completePayload))
	completeReq.Header.Set("Content-Type", "application/json")
	completeResp := httptest.NewRecorder()
	r.ServeHTTP(completeResp, completeReq)

	if completeResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", completeResp.Code)
	}

	var completed model.CheckoutSession
	if err := json.Unmarshal(completeResp.Body.Bytes(), &completed); err != nil {
		t.Fatalf("unmarshal complete response: %v", err)
	}

	if completed.Status != "completed" {
		t.Fatalf("expected status completed, got %s", completed.Status)
	}
	if completed.Order == nil || completed.Order.ID == "" {
		t.Fatalf("expected order id in response")
	}
	if orderRepo.createCount != 1 {
		t.Fatalf("expected order created once")
	}
	if paymentRepo.createCount != 1 {
		t.Fatalf("expected payment created once")
	}
	if paymentRepo.items[0].PaymentMethod != "com.nowpayments" {
		t.Fatalf("expected payment method com.nowpayments")
	}
}

func TestCheckoutCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	services := &service.Services{Checkout: checkoutService}

	handler := NewCheckoutHandler(services)

	r := gin.New()
	r.POST("/ucp/v1/checkout-sessions", handler.Create)
	r.DELETE("/ucp/v1/checkout-sessions/:id", handler.Cancel)
	r.GET("/ucp/v1/checkout-sessions/:id", handler.Get)

	createBody := model.CheckoutCreateRequest{
		Currency: "CNY",
		LineItems: []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
	}

	payload, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("marshal create request: %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions", bytes.NewReader(payload))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createResp.Code)
	}

	var created model.CheckoutSession
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	cancelReq := httptest.NewRequest(http.MethodDelete, "/ucp/v1/checkout-sessions/"+created.ID, nil)
	cancelResp := httptest.NewRecorder()
	r.ServeHTTP(cancelResp, cancelReq)

	if cancelResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", cancelResp.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/ucp/v1/checkout-sessions/"+created.ID, nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getResp.Code)
	}

	var fetched model.CheckoutSession
	if err := json.Unmarshal(getResp.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("unmarshal get response: %v", err)
	}

	if fetched.Status != "canceled" {
		t.Fatalf("expected status canceled, got %s", fetched.Status)
	}
}
