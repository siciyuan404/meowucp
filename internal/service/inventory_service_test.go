package service

import (
	"errors"
	"testing"

	"github.com/meowucp/internal/domain"
)

type atomicProductRepo struct {
	product           *domain.Product
	atomicCalled      bool
	atomicErr         error
	updatedStock      int
	updateStockCalled bool
}

func (r *atomicProductRepo) Create(product *domain.Product) error { return nil }
func (r *atomicProductRepo) Update(product *domain.Product) error { return nil }
func (r *atomicProductRepo) Delete(id int64) error                { return nil }
func (r *atomicProductRepo) FindByID(id int64) (*domain.Product, error) {
	if r.product == nil {
		return nil, errors.New("product not found")
	}
	return r.product, nil
}
func (r *atomicProductRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (r *atomicProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (r *atomicProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	if r.product == nil {
		return nil, errors.New("product not found")
	}
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}
	for _, id := range ids {
		if r.product.ID == id {
			return []*domain.Product{r.product}, nil
		}
	}
	return nil, errors.New("product not found")
}
func (r *atomicProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, nil
}
func (r *atomicProductRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *atomicProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (r *atomicProductRepo) SearchCount(query string) (int64, error) { return 0, nil }
func (r *atomicProductRepo) UpdateStock(id int64, quantity int) error {
	r.updateStockCalled = true
	r.updatedStock = quantity
	return nil
}
func (r *atomicProductRepo) IncrementViews(id int64) error               { return nil }
func (r *atomicProductRepo) IncrementSales(id int64, quantity int) error { return nil }

func (r *atomicProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	r.atomicCalled = true
	return r.atomicErr
}

type inventoryTestRepo struct {
	logs []*domain.InventoryLog
}

func (f *inventoryTestRepo) Create(log *domain.InventoryLog) error {
	f.logs = append(f.logs, log)
	return nil
}
func (f *inventoryTestRepo) FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error) {
	return nil, nil
}
func (f *inventoryTestRepo) CountByProductID(productID int64) (int64, error) { return 0, nil }

func TestInventoryServiceAdjustStockUsesAtomicUpdater(t *testing.T) {
	productRepo := &atomicProductRepo{product: &domain.Product{ID: 1, StockQuantity: 5}}
	inventoryRepo := &inventoryTestRepo{}
	svc := NewInventoryService(productRepo, inventoryRepo)

	if err := svc.AdjustStock(1, -2, "out", "ref", "order", "note"); err != nil {
		t.Fatalf("adjust stock: %v", err)
	}
	if !productRepo.atomicCalled {
		t.Fatalf("expected atomic updater to be used")
	}
	if productRepo.updateStockCalled {
		t.Fatalf("expected legacy UpdateStock not to be used")
	}
	if len(inventoryRepo.logs) != 1 {
		t.Fatalf("expected inventory log created")
	}
}

func TestInventoryServiceAdjustStockFailsWhenAtomicUpdateFails(t *testing.T) {
	productRepo := &atomicProductRepo{
		product:   &domain.Product{ID: 2, StockQuantity: 1},
		atomicErr: errors.New("insufficient stock"),
	}
	inventoryRepo := &inventoryTestRepo{}
	svc := NewInventoryService(productRepo, inventoryRepo)

	if err := svc.AdjustStock(2, -2, "out", "ref", "order", "note"); err == nil {
		t.Fatalf("expected error when atomic update fails")
	}
	if len(inventoryRepo.logs) != 0 {
		t.Fatalf("expected no inventory log on failure")
	}
}
