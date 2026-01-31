package repository

import (
	"errors"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type orderRepository struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) FindByID(id int64) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Payments").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNo(orderNo string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Payments").
		Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *orderRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	var orders []*domain.Order
	query := r.db.Model(&domain.Order{}).Order("orders.created_at DESC").Offset(offset).Limit(limit)

	if sku, ok := filters["item_sku"]; ok {
		query = query.Joins("JOIN order_items ON order_items.order_id = orders.id").Where("order_items.sku = ?", sku)
		delete(filters, "item_sku")
	}

	for key, value := range filters {
		query = query.Where(key, value)
	}

	err := query.Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Count(filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&domain.Order{})

	if sku, ok := filters["item_sku"]; ok {
		query = query.Joins("JOIN order_items ON order_items.order_id = orders.id").Where("order_items.sku = ?", sku)
		delete(filters, "item_sku")
	}

	for key, value := range filters {
		query = query.Where(key, value)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *orderRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) CreateOrderItem(item *domain.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) Transaction(fn func(orderRepo OrderRepository, cartRepo CartRepository, productRepo ProductRepository, inventoryRepo InventoryRepository, idempotencyRepo OrderIdempotencyRepository, paymentRepo PaymentRepository) error) error {
	if r.db == nil {
		return errors.New("database not initialized")
	}
	return r.db.Transaction(func(tx *database.DB) error {
		repos := NewRepositories(tx)
		return fn(repos.Order, repos.Cart, repos.Product, repos.Inventory, repos.OrderIdempotency, repos.Payment)
	})
}
