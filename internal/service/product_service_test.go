package service

import (
	"testing"

	"github.com/meowucp/internal/domain"
)

type fakeProductListRepo struct {
	listCalled int
	items      []*domain.Product
}

func (f *fakeProductListRepo) Create(product *domain.Product) error { return nil }
func (f *fakeProductListRepo) Update(product *domain.Product) error { return nil }
func (f *fakeProductListRepo) Delete(id int64) error                { return nil }
func (f *fakeProductListRepo) FindByID(id int64) (*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductListRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (f *fakeProductListRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductListRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductListRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	f.listCalled++
	return f.items, nil
}
func (f *fakeProductListRepo) Count(filters map[string]interface{}) (int64, error) {
	return int64(len(f.items)), nil
}
func (f *fakeProductListRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductListRepo) SearchCount(query string) (int64, error)  { return 0, nil }
func (f *fakeProductListRepo) UpdateStock(id int64, quantity int) error { return nil }
func (f *fakeProductListRepo) UpdateStockWithDelta(id int64, delta int) error {
	return nil
}
func (f *fakeProductListRepo) IncrementViews(id int64) error { return nil }
func (f *fakeProductListRepo) IncrementSales(id int64, quantity int) error {
	return nil
}

type fakeCache struct {
	store map[string]string
}

func newFakeCache() *fakeCache {
	return &fakeCache{store: map[string]string{}}
}

func (f *fakeCache) Get(key string) (string, bool) {
	val, ok := f.store[key]
	return val, ok
}

func (f *fakeCache) Set(key, value string) {
	f.store[key] = value
}

func TestProductListCachesResults(t *testing.T) {
	repo := &fakeProductListRepo{items: []*domain.Product{{ID: 1}, {ID: 2}}}
	cache := newFakeCache()

	svc := NewProductServiceWithCache(repo, nil, cache)
	_, _, err := svc.ListProducts(0, 10, map[string]interface{}{"status": "active"})
	if err != nil {
		t.Fatalf("list products: %v", err)
	}
	_, _, err = svc.ListProducts(0, 10, map[string]interface{}{"status": "active"})
	if err != nil {
		t.Fatalf("list products: %v", err)
	}
	if repo.listCalled != 1 {
		t.Fatalf("expected repo list called once, got %d", repo.listCalled)
	}
}
