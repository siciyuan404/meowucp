package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeLocalizationService struct{}

func (f *fakeLocalizationService) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}
	if base == "CNY" && target == "USD" {
		return amount / 7.2, nil
	}
	if base == "CNY" && target == "EUR" {
		return amount / 7.8, nil
	}
	return amount, nil
}

func (f *fakeLocalizationService) Translate(key, locale string) (string, error) {
	if locale == "en-US" && key == "product" {
		return "Product", nil
	}
	return key, nil
}

type fakePublicProductService struct {
	products []*domain.Product
}

func newFakePublicProductService() *fakePublicProductService {
	now := time.Now()
	return &fakePublicProductService{
		products: []*domain.Product{
			{ID: 1, Name: "Product 1", Slug: "product-1", Price: 720, ComparePrice: 800, SKU: "SKU-001", StockQuantity: 10, Status: 1, CreatedAt: now},
			{ID: 2, Name: "Product 2", Slug: "product-2", Price: 1440, ComparePrice: 1600, SKU: "SKU-002", StockQuantity: 20, Status: 1, CreatedAt: now},
		},
	}
}

func (f *fakePublicProductService) ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error) {
	start := offset
	if start >= len(f.products) {
		return []*domain.Product{}, 0, nil
	}
	end := offset + limit
	if end > len(f.products) {
		end = len(f.products)
	}
	return f.products[start:end], int64(len(f.products)), nil
}

func (f *fakePublicProductService) GetProduct(id int64) (*domain.Product, error) {
	for _, p := range f.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func TestProductListRespectsCurrency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	productService := newFakePublicProductService()
	localizationService := &fakeLocalizationService{}
	handler := NewProductHandler(productService, localizationService)

	r := gin.New()
	r.GET("/api/v1/products", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?currency=USD&locale=en-US", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !containsString(resp.Body.String(), "currency") {
		t.Fatalf("expected currency in response")
	}
}
