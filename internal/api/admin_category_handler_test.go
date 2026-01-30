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

type fakeCategoryService struct {
	items  map[int64]*domain.Category
	lastID int64
}

func newFakeCategoryService() *fakeCategoryService {
	return &fakeCategoryService{items: map[int64]*domain.Category{}}
}

func (f *fakeCategoryService) CreateCategory(category *domain.Category) error {
	f.lastID++
	category.ID = f.lastID
	f.items[category.ID] = category
	return nil
}

func (f *fakeCategoryService) UpdateCategory(category *domain.Category) error {
	f.items[category.ID] = category
	return nil
}

func (f *fakeCategoryService) GetCategory(id int64) (*domain.Category, error) {
	item, ok := f.items[id]
	if !ok {
		return nil, errNotFound
	}
	return item, nil
}

func (f *fakeCategoryService) ListCategories(offset, limit int) ([]*domain.Category, int64, error) {
	items := make([]*domain.Category, 0, len(f.items))
	for _, item := range f.items {
		items = append(items, item)
	}
	return items, int64(len(items)), nil
}

func TestAdminCategoryCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeCategoryService()
	handler := NewAdminCategoryHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/categories", handler.Create)

	payload := map[string]interface{}{
		"name":   "Category A",
		"slug":   "category-a",
		"status": 1,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "Category A") {
		t.Fatalf("expected response to include category name")
	}
}

func TestAdminCategoryUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeCategoryService()
	svc.CreateCategory(&domain.Category{Name: "Category A", Slug: "category-a", Status: 1})

	handler := NewAdminCategoryHandler(svc)

	r := gin.New()
	r.PUT("/api/v1/admin/categories/:id", handler.Update)

	payload := map[string]interface{}{"name": "Category B"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/categories/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if svc.items[1].Name != "Category B" {
		t.Fatalf("expected category name updated")
	}
}

func TestAdminCategoryList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := newFakeCategoryService()
	svc.CreateCategory(&domain.Category{Name: "Category A", Slug: "category-a", Status: 1})
	svc.CreateCategory(&domain.Category{Name: "Category B", Slug: "category-b", Status: 1})

	handler := NewAdminCategoryHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/categories", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/categories?page=1&limit=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "pagination") {
		t.Fatalf("expected pagination in response")
	}
}
