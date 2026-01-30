package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeProductService struct {
	items  map[int64]*domain.Product
	lastID int64
}

func newFakeProductService() *fakeProductService {
	return &fakeProductService{items: map[int64]*domain.Product{}}
}

func (f *fakeProductService) CreateProduct(product *domain.Product) error {
	f.lastID++
	product.ID = f.lastID
	f.items[product.ID] = product
	return nil
}

func (f *fakeProductService) GetProduct(id int64) (*domain.Product, error) {
	item, ok := f.items[id]
	if !ok {
		return nil, errNotFound
	}
	return item, nil
}

func (f *fakeProductService) UpdateProduct(product *domain.Product) error {
	f.items[product.ID] = product
	return nil
}

func (f *fakeProductService) ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error) {
	items := make([]*domain.Product, 0, len(f.items))
	for _, item := range f.items {
		items = append(items, item)
	}
	return items, int64(len(items)), nil
}

func TestAdminProductCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/products", handler.Create)

	payload := map[string]interface{}{
		"name":           "Test Product",
		"slug":           "test-product",
		"price":          99.5,
		"sku":            "SKU-001",
		"stock_quantity": 10,
		"status":         1,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "Test Product") {
		t.Fatalf("expected response to include product name")
	}
}

func TestAdminProductCreateErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/products", handler.Create)

	payload := map[string]interface{}{"name": ""}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.Code)
	}

	var decoded map[string]map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if decoded["error"]["code"] != "missing_required_fields" {
		t.Fatalf("expected error code missing_required_fields")
	}
}

func TestAdminProductList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	svc.CreateProduct(&domain.Product{Name: "P1", Slug: "p1", SKU: "sku1", Price: 10})
	svc.CreateProduct(&domain.Product{Name: "P2", Slug: "p2", SKU: "sku2", Price: 20})

	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/products", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products?page=1&limit=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	body := resp.Body.String()
	if !containsAll(body, []string{"pagination", "total"}) {
		t.Fatalf("expected pagination in response")
	}
}

func TestAdminProductUpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	svc.CreateProduct(&domain.Product{Name: "P1", Slug: "p1", SKU: "sku1", Price: 10, Status: 1})

	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.PATCH("/api/v1/admin/products/:id/status", handler.UpdateStatus)

	payload := map[string]interface{}{"status": 2}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/products/1/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.items[1].Status != 2 {
		t.Fatalf("expected status to be updated")
	}
}

func TestAdminProductGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	svc.CreateProduct(&domain.Product{Name: "P1", Slug: "p1", SKU: "sku1", Price: 10})

	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/products/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "P1") {
		t.Fatalf("expected response to include product")
	}
}

func TestAdminProductUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeProductService()
	svc.CreateProduct(&domain.Product{Name: "P1", Slug: "p1", SKU: "sku1", Price: 10})

	handler := NewAdminProductHandler(svc)

	r := gin.New()
	r.PUT("/api/v1/admin/products/:id", handler.Update)

	payload := map[string]interface{}{"name": "P1-Updated"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.items[1].Name != "P1-Updated" {
		t.Fatalf("expected product name updated")
	}
}

var errNotFound = gin.Error{Type: gin.ErrorTypePublic}
