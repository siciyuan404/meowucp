package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type shipmentRepository struct {
	db *database.DB
}

func NewShipmentRepository(db *database.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(shipment *domain.Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *shipmentRepository) FindByOrderID(orderID int64) (*domain.Shipment, error) {
	var shipment domain.Shipment
	if err := r.db.Where("order_id = ?", orderID).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) Update(shipment *domain.Shipment) error {
	return r.db.Save(shipment).Error
}
