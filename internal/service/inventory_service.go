package service

import (
	"errors"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type InventoryService struct {
	productRepo   repository.ProductRepository
	inventoryRepo repository.InventoryRepository
}

type atomicStockUpdater interface {
	UpdateStockWithDelta(id int64, delta int) error
}

func NewInventoryService(productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) AdjustStock(productID int64, quantity int, typeName, referenceID, referenceType, notes string) error {
	if updater, ok := s.productRepo.(atomicStockUpdater); ok {
		if err := updater.UpdateStockWithDelta(productID, quantity); err != nil {
			return err
		}
	} else {
		product, err := s.productRepo.FindByID(productID)
		if err != nil {
			return err
		}
		newStock := product.StockQuantity + quantity
		if newStock < 0 {
			return errors.New("insufficient stock")
		}

		if err := s.productRepo.UpdateStock(productID, newStock); err != nil {
			return err
		}
	}

	log := &domain.InventoryLog{
		ProductID:      productID,
		QuantityChange: quantity,
		Type:           typeName,
		ReferenceID:    referenceID,
		ReferenceType:  referenceType,
		Notes:          notes,
	}

	return s.inventoryRepo.Create(log)
}

func (s *InventoryService) GetInventoryLogs(productID int64, offset, limit int) ([]*domain.InventoryLog, int64, error) {
	logs, err := s.inventoryRepo.FindByProductID(productID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.inventoryRepo.CountByProductID(productID)
	if err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}
