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
	"github.com/jinzhu/gorm"
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

type fakeCheckoutProductRepo struct {
	products map[string]*domain.Product
	byID     map[int64]*domain.Product
}

func newFakeCheckoutProductRepo(products map[string]*domain.Product) *fakeCheckoutProductRepo {
	byID := map[int64]*domain.Product{}
	for _, product := range products {
		byID[product.ID] = product
	}
	return &fakeCheckoutProductRepo{products: products, byID: byID}
}

func (f *fakeCheckoutProductRepo) Create(product *domain.Product) error { return nil }
func (f *fakeCheckoutProductRepo) Update(product *domain.Product) error { return nil }
func (f *fakeCheckoutProductRepo) Delete(id int64) error                { return nil }
func (f *fakeCheckoutProductRepo) FindByID(id int64) (*domain.Product, error) {
	product, ok := f.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return product, nil
}
func (f *fakeCheckoutProductRepo) FindBySKU(sku string) (*domain.Product, error) {
	product, ok := f.products[sku]
	if !ok {
		return nil, errors.New("not found")
	}
	return product, nil
}
func (f *fakeCheckoutProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) Count(filters map[string]interface{}) (int64, error) {
	return 0, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) SearchCount(query string) (int64, error) {
	return 0, errors.New("not implemented")
}
func (f *fakeCheckoutProductRepo) UpdateStock(id int64, quantity int) error {
	product, ok := f.byID[id]
	if !ok {
		return errors.New("not found")
	}
	product.StockQuantity = quantity
	return nil
}
func (f *fakeCheckoutProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	product, ok := f.byID[id]
	if !ok {
		return errors.New("not found")
	}
	product.StockQuantity += delta
	if product.StockQuantity < 0 {
		return errors.New("insufficient stock")
	}
	return nil
}
func (f *fakeCheckoutProductRepo) IncrementViews(id int64) error { return nil }
func (f *fakeCheckoutProductRepo) IncrementSales(id int64, quantity int) error {
	return nil
}

type fakeCheckoutInventoryRepo struct {
	logs []*domain.InventoryLog
}

func (f *fakeCheckoutInventoryRepo) Create(log *domain.InventoryLog) error {
	f.logs = append(f.logs, log)
	return nil
}
func (f *fakeCheckoutInventoryRepo) FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error) {
	return []*domain.InventoryLog{}, nil
}
func (f *fakeCheckoutInventoryRepo) CountByProductID(productID int64) (int64, error) {
	return int64(len(f.logs)), nil
}

type fakeCheckoutIdempotencyRepo struct {
	records map[string]*domain.OrderIdempotency
}

func newFakeCheckoutIdempotencyRepo() *fakeCheckoutIdempotencyRepo {
	return &fakeCheckoutIdempotencyRepo{records: map[string]*domain.OrderIdempotency{}}
}

func (f *fakeCheckoutIdempotencyRepo) Create(record *domain.OrderIdempotency) error {
	key := record.IdempotencyKey
	if _, exists := f.records[key]; exists {
		return errors.New("duplicate")
	}
	f.records[key] = record
	return nil
}
func (f *fakeCheckoutIdempotencyRepo) FindByUserIDAndIdempotencyKey(userID int64, key string) (*domain.OrderIdempotency, error) {
	record, ok := f.records[key]
	if !ok || record.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return record, nil
}
func (f *fakeCheckoutIdempotencyRepo) Update(record *domain.OrderIdempotency) error {
	f.records[record.IdempotencyKey] = record
	return nil
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
	if created.ContinueURL != "" {
		t.Fatalf("expected continue_url to be empty for ready_for_complete")
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

	paymentHandlersMissing := false
	for _, message := range created.Messages {
		if message.Code == "payment_required" {
			t.Fatalf("expected payment_required message to be omitted")
		}
		if message.Code == "payment_handlers_missing" && message.Severity == "requires_buyer_input" {
			paymentHandlersMissing = true
		}
	}
	if !paymentHandlersMissing {
		t.Fatalf("expected payment_handlers_missing requires_buyer_input message")
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

func TestCheckoutUpdateRequiresEscalationWhenSignInNeeded(t *testing.T) {
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

	updateBody := map[string]any{
		"id":       created.ID,
		"currency": "CNY",
		"line_items": []model.LineItem{
			{
				Item: model.Item{
					ID:    "sku_1",
					Title: "Test Item",
					Price: 19900,
				},
				Quantity: 1,
			},
		},
		"requires_sign_in": true,
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

	if updated.Status != "requires_escalation" {
		t.Fatalf("expected status requires_escalation, got %s", updated.Status)
	}

	requiresSignIn := false
	for _, message := range updated.Messages {
		if message.Code == "requires_sign_in" && message.Severity == "requires_buyer_input" {
			requiresSignIn = true
			break
		}
	}
	if !requiresSignIn {
		t.Fatalf("expected requires_sign_in requires_buyer_input message")
	}
}

func TestCheckoutCompleteCreatesOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	orderRepo := newFakeOrderRepo()
	productRepo := newFakeCheckoutProductRepo(map[string]*domain.Product{
		"sku_1": {ID: 10, Name: "Test Item", SKU: "sku_1", StockQuantity: 5},
	})
	inventoryRepo := &fakeCheckoutInventoryRepo{}
	idempotencyRepo := newFakeCheckoutIdempotencyRepo()

	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	orderService := service.NewOrderService(orderRepo, nil, productRepo, inventoryRepo, idempotencyRepo)
	services := &service.Services{
		Checkout: checkoutService,
		Order:    orderService,
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

	if completed.Order == nil || completed.Order.ID == "" {
		t.Fatalf("expected order id in response")
	}
	if orderRepo.createCount != 1 {
		t.Fatalf("expected order created once")
	}
}

func TestCheckoutCompleteCreatesOrderAndPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	orderRepo := newFakeOrderRepo()
	paymentRepo := newFakePaymentRepo()
	productRepo := newFakeCheckoutProductRepo(map[string]*domain.Product{
		"sku_1": {ID: 10, Name: "Test Item", SKU: "sku_1", StockQuantity: 5},
	})
	inventoryRepo := &fakeCheckoutInventoryRepo{}
	idempotencyRepo := newFakeCheckoutIdempotencyRepo()
	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	orderService := service.NewOrderService(orderRepo, nil, productRepo, inventoryRepo, idempotencyRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)

	services := &service.Services{
		Checkout: checkoutService,
		Order:    orderService,
		Payment:  paymentService,
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

func TestCheckoutCompleteIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	checkoutRepo := newFakeCheckoutRepo()
	orderRepo := newFakeOrderRepo()
	productRepo := newFakeCheckoutProductRepo(map[string]*domain.Product{
		"sku_1": {ID: 10, Name: "Test Item", SKU: "sku_1", StockQuantity: 5},
	})
	inventoryRepo := &fakeCheckoutInventoryRepo{}
	idempotencyRepo := newFakeCheckoutIdempotencyRepo()

	checkoutService := service.NewCheckoutSessionService(checkoutRepo)
	orderService := service.NewOrderService(orderRepo, nil, productRepo, inventoryRepo, idempotencyRepo)
	services := &service.Services{
		Checkout: checkoutService,
		Order:    orderService,
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

	var first model.CheckoutSession
	if err := json.Unmarshal(completeResp.Body.Bytes(), &first); err != nil {
		t.Fatalf("unmarshal first complete response: %v", err)
	}

	completeReq2 := httptest.NewRequest(http.MethodPost, "/ucp/v1/checkout-sessions/"+created.ID+"/complete", bytes.NewReader(completePayload))
	completeReq2.Header.Set("Content-Type", "application/json")
	completeResp2 := httptest.NewRecorder()
	r.ServeHTTP(completeResp2, completeReq2)

	if completeResp2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", completeResp2.Code)
	}

	var second model.CheckoutSession
	if err := json.Unmarshal(completeResp2.Body.Bytes(), &second); err != nil {
		t.Fatalf("unmarshal second complete response: %v", err)
	}

	if first.Order == nil || second.Order == nil || first.Order.ID == "" {
		t.Fatalf("expected order ids to be set")
	}
	if first.Order.ID != second.Order.ID {
		t.Fatalf("expected same order id on retry")
	}
	if orderRepo.createCount != 1 {
		t.Fatalf("expected order created once")
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
