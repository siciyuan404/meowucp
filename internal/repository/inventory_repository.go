package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type inventoryRepository struct {
	db *database.DB
}

func NewInventoryRepository(db *database.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(log *domain.InventoryLog) error {
	return r.db.Create(log).Error
}

func (r *inventoryRepository) FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error) {
	var logs []*domain.InventoryLog
	err := r.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *inventoryRepository) CountByProductID(productID int64) (int64, error) {
	var count int64
	err := r.db.Model(&domain.InventoryLog{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}
