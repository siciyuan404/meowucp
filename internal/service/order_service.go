package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type OrderService struct {
	orderRepo       repository.OrderRepository
	cartRepo        repository.CartRepository
	productRepo     repository.ProductRepository
	inventoryRepo   repository.InventoryRepository
	idempotencyRepo repository.OrderIdempotencyRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, idempotencyRepo repository.OrderIdempotencyRepository) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		cartRepo:        cartRepo,
		productRepo:     productRepo,
		inventoryRepo:   inventoryRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

type orderTransactionRunner interface {
	Transaction(fn func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, idempotencyRepo repository.OrderIdempotencyRepository, paymentRepo repository.PaymentRepository) error) error
}

var ErrOrderIdempotencyConflict = errors.New("order idempotency conflict")

func (s *OrderService) CreateOrder(userID int64, idempotencyKey string, shippingAddress, billingAddress, notes string, paymentMethod string) (*domain.Order, error) {
	if txRunner, ok := s.orderRepo.(orderTransactionRunner); ok {
		var createdOrder *domain.Order
		err := txRunner.Transaction(func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, idempotencyRepo repository.OrderIdempotencyRepository, paymentRepo repository.PaymentRepository) error {
			var err error
			createdOrder, err = s.createOrderWithRepos(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo, userID, idempotencyKey, shippingAddress, billingAddress, notes, paymentMethod)
			return err
		})
		if err != nil {
			return nil, err
		}
		return createdOrder, nil
	}

	return s.createOrderWithRepos(s.orderRepo, s.cartRepo, s.productRepo, s.inventoryRepo, s.idempotencyRepo, userID, idempotencyKey, shippingAddress, billingAddress, notes, paymentMethod)
}

func (s *OrderService) createOrderWithRepos(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, idempotencyRepo repository.OrderIdempotencyRepository, userID int64, idempotencyKey string, shippingAddress, billingAddress, notes string, paymentMethod string) (*domain.Order, error) {
	var idempotencyRecord *domain.OrderIdempotency
	if idempotencyKey != "" {
		if idempotencyRepo == nil {
			return nil, errors.New("idempotency repository unavailable")
		}
		record, err := idempotencyRepo.FindByUserIDAndIdempotencyKey(userID, idempotencyKey)
		if err == nil {
			if record.OrderID != nil {
				order, err := orderRepo.FindByID(*record.OrderID)
				if err != nil {
					return nil, err
				}
				return order, nil
			}
			return nil, ErrOrderIdempotencyConflict
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		record = &domain.OrderIdempotency{
			UserID:         userID,
			IdempotencyKey: idempotencyKey,
			Status:         "pending",
		}
		if err := idempotencyRepo.Create(record); err != nil {
			return nil, err
		}
		idempotencyRecord = record
	}

	if cartRepo == nil || orderRepo == nil || productRepo == nil || inventoryRepo == nil {
		return nil, errors.New("order dependencies unavailable")
	}
	cart, err := cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	subtotal := 0.0
	orderItems := make([]domain.OrderItem, 0, len(cart.Items))
	productIDs := make([]int64, 0, len(cart.Items))
	productIDSet := map[int64]struct{}{}
	for _, item := range cart.Items {
		if _, exists := productIDSet[item.ProductID]; exists {
			continue
		}
		productIDSet[item.ProductID] = struct{}{}
		productIDs = append(productIDs, item.ProductID)
	}
	products, err := productRepo.GetByIDs(productIDs)
	if err != nil {
		return nil, errors.New("product not found")
	}
	productByID := make(map[int64]*domain.Product, len(products))
	for _, product := range products {
		productByID[product.ID] = product
	}

	for _, item := range cart.Items {
		product, ok := productByID[item.ProductID]
		if !ok {
			return nil, errors.New("product not found")
		}
		productID := item.ProductID

		subtotal += item.Price * float64(item.Quantity)
		orderItems = append(orderItems, domain.OrderItem{
			ProductID:   &productID,
			ProductName: product.Name,
			SKU:         product.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.Price,
			TotalPrice:  item.Price * float64(item.Quantity),
		})
	}

	shippingFee := 10.0
	tax := subtotal * 0.1
	total := subtotal + shippingFee + tax

	orderNo := fmt.Sprintf("ORD%d%04d", time.Now().Unix(), userID%10000)

	order := &domain.Order{
		OrderNo:         orderNo,
		UserID:          &userID,
		Status:          "pending",
		Subtotal:        subtotal,
		ShippingFee:     shippingFee,
		Tax:             tax,
		Total:           total,
		PaymentMethod:   paymentMethod,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		Notes:           notes,
	}

	if err := orderRepo.Create(order); err != nil {
		return nil, err
	}

	inventorySvc := NewInventoryService(productRepo, inventoryRepo)
	for _, item := range orderItems {
		item.OrderID = order.ID
		if err := orderRepo.CreateOrderItem(&item); err != nil {
			return nil, err
		}
		if err := inventorySvc.AdjustStock(
			*item.ProductID,
			-item.Quantity,
			"out",
			orderNo,
			"order",
			fmt.Sprintf("Order %s created", orderNo),
		); err != nil {
			return nil, err
		}
		if err := productRepo.IncrementSales(*item.ProductID, item.Quantity); err != nil {
			return nil, err
		}
	}

	if err := cartRepo.ClearCart(cart.ID); err != nil {
		return nil, err
	}

	if idempotencyRecord != nil {
		orderID := order.ID
		idempotencyRecord.OrderID = &orderID
		idempotencyRecord.Status = "completed"
		if err := idempotencyRepo.Update(idempotencyRecord); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (s *OrderService) GetOrder(id int64) (*domain.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) GetOrderByOrderNo(orderNo string) (*domain.Order, error) {
	return s.orderRepo.FindByOrderNo(orderNo)
}

func (s *OrderService) ListUserOrders(userID int64, offset, limit int) ([]*domain.Order, int64, error) {
	orders, err := s.orderRepo.FindByUserID(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.orderRepo.CountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (s *OrderService) UpdateOrderStatus(id int64, status string) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) ListOrders(offset, limit int, filters map[string]interface{}) ([]*domain.Order, int64, error) {
	filterCopy := map[string]interface{}{}
	for key, value := range filters {
		filterCopy[key] = value
	}
	orders, err := s.orderRepo.List(offset, limit, filterCopy)
	if err != nil {
		return nil, 0, err
	}
	countFilters := map[string]interface{}{}
	for key, value := range filters {
		countFilters[key] = value
	}
	count, err := s.orderRepo.Count(countFilters)
	if err != nil {
		return nil, 0, err
	}
	return orders, count, nil
}
