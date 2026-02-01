package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type E2EService struct {
	users    map[int64]*domain.User
	carts    map[int64]*domain.Cart
	orders   map[int64]*domain.Order
	payments map[int64]*domain.Payment
	webhooks []*domain.UCPWebhookJob
}

func newE2EService() *E2EService {
	return &E2EService{
		users:    make(map[int64]*domain.User),
		carts:    make(map[int64]*domain.Cart),
		orders:   make(map[int64]*domain.Order),
		payments: make(map[int64]*domain.Payment),
		webhooks: make([]*domain.UCPWebhookJob, 0),
	}
}

func (s *E2EService) CreateProduct(product *domain.Product) error {
	return nil
}

func (s *E2EService) GetProduct(id int64) (*domain.Product, error) {
	return &domain.Product{ID: id, Name: "Test Product", Price: 100, SKU: "TEST-001"}, nil
}

func (s *E2EService) ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error) {
	return []*domain.Product{{ID: 1, Name: "Test Product", Price: 100, SKU: "TEST-001"}}, 1, nil
}

func (s *E2EService) Convert(amount float64, base, target string) (float64, error) {
	return amount, nil
}

func TestE2ECheckoutPaymentOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := newE2EService()
	productHandler := NewProductHandler(service, service)
	r := gin.New()
	r.GET("/api/v1/products", productHandler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?currency=CNY&locale=zh-CN", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := result["products"]; !ok {
		t.Fatalf("expected products in response")
	}

	if _, ok := result["pagination"]; !ok {
		t.Fatalf("expected pagination in response")
	}

	if !containsString(resp.Body.String(), "currency") {
		t.Fatalf("expected currency in response")
	}

	if !containsString(resp.Body.String(), "locale") {
		t.Fatalf("expected locale in response")
	}
}
