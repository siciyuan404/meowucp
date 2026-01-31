package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/model"
)

type fakeOrderRepo struct {
	items []*domain.Order
}

func (f *fakeOrderRepo) Create(order *domain.Order) error                    { return nil }
func (f *fakeOrderRepo) Update(order *domain.Order) error                    { return nil }
func (f *fakeOrderRepo) FindByID(id int64) (*domain.Order, error)            { return nil, nil }
func (f *fakeOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (f *fakeOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakeOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	if offset >= len(f.items) {
		return []*domain.Order{}, nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], nil
}
func (f *fakeOrderRepo) Count(filters map[string]interface{}) (int64, error) {
	return int64(len(f.items)), nil
}
func (f *fakeOrderRepo) UpdateStatus(id int64, status string) error   { return nil }
func (f *fakeOrderRepo) CreateOrderItem(item *domain.OrderItem) error { return nil }

type fakePaidOrderRepo struct {
	order          *domain.Order
	updateCalled   bool
	updateStatusID int64
	updateStatus   string
}

func (f *fakePaidOrderRepo) Create(order *domain.Order) error { return nil }
func (f *fakePaidOrderRepo) Update(order *domain.Order) error {
	f.order = order
	f.updateCalled = true
	return nil
}
func (f *fakePaidOrderRepo) FindByID(id int64) (*domain.Order, error)            { return f.order, nil }
func (f *fakePaidOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (f *fakePaidOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakePaidOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakePaidOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakePaidOrderRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (f *fakePaidOrderRepo) UpdateStatus(id int64, status string) error {
	f.updateStatusID = id
	f.updateStatus = status
	return nil
}
func (f *fakePaidOrderRepo) CreateOrderItem(item *domain.OrderItem) error { return nil }

type fakeWebhookQueueRepo struct {
	last *domain.UCPWebhookJob
}

func (f *fakeWebhookQueueRepo) Create(job *domain.UCPWebhookJob) error { f.last = job; return nil }
func (f *fakeWebhookQueueRepo) ListDue(limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}
func (f *fakeWebhookQueueRepo) Update(job *domain.UCPWebhookJob) error { f.last = job; return nil }
func (f *fakeWebhookQueueRepo) List(offset, limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}
func (f *fakeWebhookQueueRepo) Count() (int64, error) { return 0, nil }
func (f *fakeWebhookQueueRepo) FindByID(id int64) (*domain.UCPWebhookJob, error) {
	return nil, errors.New("not found")
}

type fakeOrderCreateRepo struct {
	createdOrder *domain.Order
	createdItems []*domain.OrderItem
	createErr    error
	itemErr      error
}

func (f *fakeOrderCreateRepo) Create(order *domain.Order) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createdOrder = order
	return nil
}
func (f *fakeOrderCreateRepo) Update(order *domain.Order) error                    { return nil }
func (f *fakeOrderCreateRepo) FindByID(id int64) (*domain.Order, error)            { return nil, nil }
func (f *fakeOrderCreateRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (f *fakeOrderCreateRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeOrderCreateRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakeOrderCreateRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeOrderCreateRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (f *fakeOrderCreateRepo) UpdateStatus(id int64, status string) error          { return nil }
func (f *fakeOrderCreateRepo) CreateOrderItem(item *domain.OrderItem) error {
	if f.itemErr != nil {
		return f.itemErr
	}
	f.createdItems = append(f.createdItems, item)
	return nil
}

type fakeIdempotencyOrderRepo struct {
	order        *domain.Order
	createCalled bool
}

func (f *fakeIdempotencyOrderRepo) Create(order *domain.Order) error {
	f.createCalled = true
	return errors.New("unexpected create")
}
func (f *fakeIdempotencyOrderRepo) Update(order *domain.Order) error         { return nil }
func (f *fakeIdempotencyOrderRepo) FindByID(id int64) (*domain.Order, error) { return f.order, nil }
func (f *fakeIdempotencyOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) {
	return nil, nil
}
func (f *fakeIdempotencyOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeIdempotencyOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakeIdempotencyOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeIdempotencyOrderRepo) Count(filters map[string]interface{}) (int64, error) {
	return 0, nil
}
func (f *fakeIdempotencyOrderRepo) UpdateStatus(id int64, status string) error   { return nil }
func (f *fakeIdempotencyOrderRepo) CreateOrderItem(item *domain.OrderItem) error { return nil }

type fakeOrderStatusRepo struct {
	order *domain.Order
}

func (f *fakeOrderStatusRepo) Create(order *domain.Order) error                    { return nil }
func (f *fakeOrderStatusRepo) Update(order *domain.Order) error                    { f.order = order; return nil }
func (f *fakeOrderStatusRepo) FindByID(id int64) (*domain.Order, error)            { return f.order, nil }
func (f *fakeOrderStatusRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (f *fakeOrderStatusRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeOrderStatusRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (f *fakeOrderStatusRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (f *fakeOrderStatusRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (f *fakeOrderStatusRepo) UpdateStatus(id int64, status string) error {
	if f.order != nil {
		f.order.Status = status
	}
	return nil
}
func (f *fakeOrderStatusRepo) CreateOrderItem(item *domain.OrderItem) error { return nil }

type fakeOrderWebhookQueueRepo struct {
	last *domain.UCPWebhookJob
}

func (f *fakeOrderWebhookQueueRepo) Create(job *domain.UCPWebhookJob) error { f.last = job; return nil }
func (f *fakeOrderWebhookQueueRepo) ListDue(limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}
func (f *fakeOrderWebhookQueueRepo) Update(job *domain.UCPWebhookJob) error { f.last = job; return nil }
func (f *fakeOrderWebhookQueueRepo) List(offset, limit int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}
func (f *fakeOrderWebhookQueueRepo) Count() (int64, error) { return 0, nil }
func (f *fakeOrderWebhookQueueRepo) FindByID(id int64) (*domain.UCPWebhookJob, error) {
	return nil, errors.New("not found")
}

type fakeOrderIdempotencyRepo struct {
	record     *domain.OrderIdempotency
	createErr  error
	updateErr  error
	createCall int
	updateCall int
}

func (f *fakeOrderIdempotencyRepo) Create(record *domain.OrderIdempotency) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.record = record
	f.createCall++
	return nil
}

func (f *fakeOrderIdempotencyRepo) FindByUserIDAndIdempotencyKey(userID int64, key string) (*domain.OrderIdempotency, error) {
	if f.record == nil || f.record.UserID != userID || f.record.IdempotencyKey != key {
		return nil, gorm.ErrRecordNotFound
	}
	return f.record, nil
}

func (f *fakeOrderIdempotencyRepo) Update(record *domain.OrderIdempotency) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.record = record
	f.updateCall++
	return nil
}

type raceOrderIdempotencyRepo struct {
	record      *domain.OrderIdempotency
	createErr   error
	createCalls int
	findCalls   int
}

func (r *raceOrderIdempotencyRepo) Create(record *domain.OrderIdempotency) error {
	r.createCalls++
	return r.createErr
}

func (r *raceOrderIdempotencyRepo) FindByUserIDAndIdempotencyKey(userID int64, key string) (*domain.OrderIdempotency, error) {
	r.findCalls++
	if r.findCalls == 1 {
		return nil, gorm.ErrRecordNotFound
	}
	if r.record == nil || r.record.UserID != userID || r.record.IdempotencyKey != key {
		return nil, gorm.ErrRecordNotFound
	}
	return r.record, nil
}

func (r *raceOrderIdempotencyRepo) Update(record *domain.OrderIdempotency) error {
	r.record = record
	return nil
}

type fakeCartRepo struct {
	cart        *domain.Cart
	clearedCart int64
	clearErr    error
}

func (f *fakeCartRepo) FindByUserID(userID int64) (*domain.Cart, error) { return f.cart, nil }
func (f *fakeCartRepo) Create(cart *domain.Cart) error                  { return nil }
func (f *fakeCartRepo) Update(cart *domain.Cart) error                  { return nil }
func (f *fakeCartRepo) Delete(id int64) error                           { return nil }
func (f *fakeCartRepo) AddItem(item *domain.CartItem) error             { return nil }
func (f *fakeCartRepo) UpdateItem(item *domain.CartItem) error          { return nil }
func (f *fakeCartRepo) RemoveItem(cartID, productID int64) error        { return nil }
func (f *fakeCartRepo) ClearCart(cartID int64) error {
	if f.clearErr != nil {
		return f.clearErr
	}
	f.clearedCart = cartID
	return nil
}

type fakeProductRepo struct {
	products       map[int64]*domain.Product
	updateStockErr error
	updatedStock   map[int64]int
	sales          map[int64]int
}

func (f *fakeProductRepo) Create(product *domain.Product) error { return nil }
func (f *fakeProductRepo) Update(product *domain.Product) error { return nil }
func (f *fakeProductRepo) Delete(id int64) error                { return nil }
func (f *fakeProductRepo) FindByID(id int64) (*domain.Product, error) {
	product, ok := f.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}
func (f *fakeProductRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (f *fakeProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0, len(ids))
	for _, id := range ids {
		product, ok := f.products[id]
		if !ok {
			return nil, errors.New("product not found")
		}
		products = append(products, product)
	}
	return products, nil
}
func (f *fakeProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (f *fakeProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (f *fakeProductRepo) SearchCount(query string) (int64, error) { return 0, nil }
func (f *fakeProductRepo) UpdateStock(id int64, quantity int) error {
	if f.updateStockErr != nil {
		return f.updateStockErr
	}
	if f.updatedStock == nil {
		f.updatedStock = map[int64]int{}
	}
	f.updatedStock[id] = quantity
	return nil
}
func (f *fakeProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	product, ok := f.products[id]
	if !ok {
		return errors.New("product not found")
	}
	newStock := product.StockQuantity + delta
	if newStock < 0 {
		return errors.New("insufficient stock")
	}
	product.StockQuantity = newStock
	return f.UpdateStock(id, newStock)
}
func (f *fakeProductRepo) IncrementViews(id int64) error { return nil }
func (f *fakeProductRepo) IncrementSales(id int64, quantity int) error {
	if f.sales == nil {
		f.sales = map[int64]int{}
	}
	f.sales[id] += quantity
	return nil
}

type fakeInventoryRepo struct {
	logs      []*domain.InventoryLog
	createErr error
}

func (f *fakeInventoryRepo) Create(log *domain.InventoryLog) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.logs = append(f.logs, log)
	return nil
}
func (f *fakeInventoryRepo) FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error) {
	return nil, nil
}
func (f *fakeInventoryRepo) CountByProductID(productID int64) (int64, error) { return 0, nil }

type batchOnlyProductRepo struct {
	products map[int64]*domain.Product
}

func (r *batchOnlyProductRepo) Create(product *domain.Product) error { return nil }
func (r *batchOnlyProductRepo) Update(product *domain.Product) error { return nil }
func (r *batchOnlyProductRepo) Delete(id int64) error                { return nil }
func (r *batchOnlyProductRepo) FindByID(id int64) (*domain.Product, error) {
	return nil, errors.New("batch lookup required")
}
func (r *batchOnlyProductRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (r *batchOnlyProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (r *batchOnlyProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, nil
}
func (r *batchOnlyProductRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *batchOnlyProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (r *batchOnlyProductRepo) SearchCount(query string) (int64, error)  { return 0, nil }
func (r *batchOnlyProductRepo) UpdateStock(id int64, quantity int) error { return nil }
func (r *batchOnlyProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	return nil
}
func (r *batchOnlyProductRepo) IncrementViews(id int64) error               { return nil }
func (r *batchOnlyProductRepo) IncrementSales(id int64, quantity int) error { return nil }
func (r *batchOnlyProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0, len(ids))
	for _, id := range ids {
		product, ok := r.products[id]
		if !ok {
			return nil, errors.New("product not found")
		}
		products = append(products, product)
	}
	return products, nil
}

func TestOrderServiceListOrders(t *testing.T) {
	repo := &fakeOrderRepo{items: []*domain.Order{{ID: 1}, {ID: 2}, {ID: 3}}}
	svc := NewOrderService(repo, nil, nil, nil, nil)

	items, total, err := svc.ListOrders(1, 1, map[string]interface{}{})
	if err != nil {
		t.Fatalf("list orders: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3")
	}
	if len(items) != 1 || items[0].ID != 2 {
		t.Fatalf("expected item id 2")
	}
}

func TestOrderServiceCreateOrderReturnsExistingOrderForIdempotencyKey(t *testing.T) {
	orderID := int64(101)
	existingOrder := &domain.Order{ID: orderID, OrderNo: "ORD-existing"}
	orderRepo := &fakeIdempotencyOrderRepo{order: existingOrder}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		10: {
			ID:            10,
			Name:          "Cat Toy",
			SKU:           "CAT-TOY-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     100,
		UserID: 1,
		Items: []domain.CartItem{{
			ProductID: 10,
			Quantity:  1,
			Price:     10,
		}},
	}}
	idempotencyRepo := &fakeOrderIdempotencyRepo{record: &domain.OrderIdempotency{
		UserID:         1,
		IdempotencyKey: "key-1",
		OrderID:        &orderID,
		Status:         "completed",
	}}
	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo)

	order, err := svc.CreateOrder(1, "key-1", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if order == nil || order.ID != orderID {
		t.Fatalf("expected existing order to be returned")
	}
	if orderRepo.createCalled {
		t.Fatalf("expected no new order creation")
	}
}

func TestOrderServiceCreateOrderReturnsConflictForPendingIdempotencyKey(t *testing.T) {
	orderRepo := &fakeIdempotencyOrderRepo{}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		20: {
			ID:            20,
			Name:          "Cat Treats",
			SKU:           "TREATS-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     200,
		UserID: 2,
		Items: []domain.CartItem{{
			ProductID: 20,
			Quantity:  1,
			Price:     12,
		}},
	}}
	idempotencyRepo := &fakeOrderIdempotencyRepo{record: &domain.OrderIdempotency{
		UserID:         2,
		IdempotencyKey: "key-2",
		Status:         "pending",
	}}
	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo)

	_, err := svc.CreateOrder(2, "key-2", "ship", "bill", "", "card")
	if !errors.Is(err, ErrOrderIdempotencyConflict) {
		t.Fatalf("expected idempotency conflict error")
	}
	if orderRepo.createCalled {
		t.Fatalf("expected no new order creation")
	}
}

func TestOrderServiceCreateOrderHandlesIdempotencyCreateCollisionWithExistingOrder(t *testing.T) {
	orderID := int64(101)
	existingOrder := &domain.Order{ID: orderID, OrderNo: "ORD-existing"}
	orderRepo := &fakeIdempotencyOrderRepo{order: existingOrder}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		10: {
			ID:            10,
			Name:          "Cat Toy",
			SKU:           "CAT-TOY-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     100,
		UserID: 1,
		Items: []domain.CartItem{{
			ProductID: 10,
			Quantity:  1,
			Price:     10,
		}},
	}}
	idempotencyRepo := &raceOrderIdempotencyRepo{
		record: &domain.OrderIdempotency{
			UserID:         1,
			IdempotencyKey: "key-dup",
			OrderID:        &orderID,
			Status:         "completed",
		},
		createErr: errors.New("duplicate key value violates unique constraint"),
	}
	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo)

	order, err := svc.CreateOrder(1, "key-dup", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if order == nil || order.ID != orderID {
		t.Fatalf("expected existing order to be returned")
	}
	if orderRepo.createCalled {
		t.Fatalf("expected no new order creation")
	}
}

func TestOrderServiceCreateOrderHandlesIdempotencyCreateCollisionWithPendingRecord(t *testing.T) {
	orderRepo := &fakeIdempotencyOrderRepo{}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		20: {
			ID:            20,
			Name:          "Cat Treats",
			SKU:           "TREATS-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     200,
		UserID: 2,
		Items: []domain.CartItem{{
			ProductID: 20,
			Quantity:  1,
			Price:     12,
		}},
	}}
	idempotencyRepo := &raceOrderIdempotencyRepo{
		record: &domain.OrderIdempotency{
			UserID:         2,
			IdempotencyKey: "key-dup-pending",
			Status:         "pending",
		},
		createErr: errors.New("duplicate key value violates unique constraint"),
	}
	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo)

	_, err := svc.CreateOrder(2, "key-dup-pending", "ship", "bill", "", "card")
	if !errors.Is(err, ErrOrderIdempotencyConflict) {
		t.Fatalf("expected idempotency conflict error")
	}
	if orderRepo.createCalled {
		t.Fatalf("expected no new order creation")
	}
}

func TestPaymentCallbackTriggersPaidWebhook(t *testing.T) {
	order := &domain.Order{ID: 42, Status: "pending"}
	orderRepo := &fakePaidOrderRepo{order: order}
	webhookRepo := &fakeWebhookQueueRepo{}
	webhookQueue := NewWebhookQueueService(webhookRepo)

	svc := NewOrderService(orderRepo, nil, nil, nil, nil)
	svc.SetWebhookQueue(webhookQueue)

	if err := svc.UpdateOrderStatus(42, "paid"); err != nil {
		t.Fatalf("update order status: %v", err)
	}
	if webhookRepo.last == nil {
		t.Fatalf("expected webhook job to be enqueued")
	}

	var payload model.OrderWebhookEvent
	if err := json.Unmarshal([]byte(webhookRepo.last.Payload), &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.EventID == "" {
		t.Fatalf("expected event_id to be set")
	}
	if payload.EventType != "order.paid" {
		t.Fatalf("expected event_type order.paid")
	}
	if payload.Order.ID != "42" || payload.Order.Status != "paid" {
		t.Fatalf("expected paid order payload")
	}
}

func TestOrderServiceEnqueuesOrderWebhookOnStatusChange(t *testing.T) {
	order := &domain.Order{ID: 1, OrderNo: "ORD-1", Status: "pending"}
	orderRepo := &fakeOrderStatusRepo{order: order}
	queueRepo := &fakeOrderWebhookQueueRepo{}
	queueService := NewWebhookQueueService(queueRepo)

	svc := NewOrderService(orderRepo, nil, nil, nil, nil)
	svc.SetWebhookQueue(queueService)

	if err := svc.UpdateOrderStatus(1, "paid"); err != nil {
		t.Fatalf("update status: %v", err)
	}
	if queueRepo.last == nil {
		t.Fatalf("expected webhook job to be enqueued")
	}

	var payload orderWebhookPayload
	if err := json.Unmarshal([]byte(queueRepo.last.Payload), &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.EventType != "paid" {
		t.Fatalf("expected event_type paid, got %s", payload.EventType)
	}
	if payload.OrderNo != "ORD-1" {
		t.Fatalf("expected order_no ORD-1, got %s", payload.OrderNo)
	}
}

func TestOrderServiceCreateOrderMarksIdempotencyFailedOnNonTransactionError(t *testing.T) {
	orderRepo := &fakeOrderCreateRepo{createErr: errors.New("order create failed")}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		10: {
			ID:            10,
			Name:          "Cat Toy",
			SKU:           "CAT-TOY-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     100,
		UserID: 1,
		Items: []domain.CartItem{{
			ProductID: 10,
			Quantity:  1,
			Price:     10,
		}},
	}}
	idempotencyRepo := &fakeOrderIdempotencyRepo{}
	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, idempotencyRepo)

	_, err := svc.CreateOrder(1, "key-fail", "ship", "bill", "", "card")
	if err == nil {
		t.Fatalf("expected order create error")
	}
	if idempotencyRepo.updateCall != 1 {
		t.Fatalf("expected idempotency record to be marked failed")
	}
	if idempotencyRepo.record == nil || idempotencyRepo.record.Status != "failed" {
		t.Fatalf("expected idempotency record status to be failed")
	}
}

func TestOrderServiceCreateOrderFillsProductInfoAndAdjustsStock(t *testing.T) {
	orderRepo := &fakeOrderCreateRepo{}
	productRepo := &fakeProductRepo{products: map[int64]*domain.Product{
		10: {
			ID:            10,
			Name:          "Cat Toy",
			SKU:           "CAT-TOY-001",
			StockQuantity: 5,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     100,
		UserID: 1,
		Items: []domain.CartItem{{
			ProductID: 10,
			Quantity:  2,
			Price:     12.5,
		}},
	}}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	order, err := svc.CreateOrder(1, "", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if order == nil || orderRepo.createdOrder == nil {
		t.Fatalf("expected order created")
	}
	if len(orderRepo.createdItems) != 1 {
		t.Fatalf("expected 1 order item")
	}
	item := orderRepo.createdItems[0]
	if item.ProductName != "Cat Toy" {
		t.Fatalf("expected product name from product repo")
	}
	if item.SKU != "CAT-TOY-001" {
		t.Fatalf("expected sku from product repo")
	}
	if productRepo.updatedStock[10] != 3 {
		t.Fatalf("expected stock to be updated to 3")
	}
	if len(inventoryRepo.logs) != 1 {
		t.Fatalf("expected inventory log to be created")
	}
	log := inventoryRepo.logs[0]
	if log.QuantityChange != -2 || log.Type != "out" {
		t.Fatalf("expected inventory log to record stock out")
	}
	if log.ReferenceID != order.OrderNo || log.ReferenceType != "order" {
		t.Fatalf("expected inventory log to reference order")
	}
	if cartRepo.clearedCart != 100 {
		t.Fatalf("expected cart to be cleared")
	}
}

func TestOrderServiceCreateOrderUsesBatchProductLookup(t *testing.T) {
	orderRepo := &fakeOrderCreateRepo{}
	productRepo := &batchOnlyProductRepo{products: map[int64]*domain.Product{
		99: {
			ID:            99,
			Name:          "Cat Tree",
			SKU:           "TREE-001",
			StockQuantity: 6,
		},
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     900,
		UserID: 9,
		Items: []domain.CartItem{{
			ProductID: 99,
			Quantity:  1,
			Price:     88,
		}},
	}}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(9, "", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
}

func TestOrderServiceCreateOrderRejectsInsufficientStock(t *testing.T) {
	store := &orderTestStore{
		products: map[int64]*domain.Product{
			20: {
				ID:            20,
				Name:          "Cat Food",
				SKU:           "FOOD-001",
				StockQuantity: 1,
			},
		},
		cart: &domain.Cart{
			ID:     200,
			UserID: 2,
			Items: []domain.CartItem{{
				ProductID: 20,
				Quantity:  2,
				Price:     10,
			}},
		},
	}
	orderRepo := &txOrderRepo{store: store}
	productRepo := &txProductRepo{store: store}
	inventoryRepo := &txInventoryRepo{store: store}
	cartRepo := &txCartRepo{store: store}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(2, "", "ship", "bill", "", "card")
	if err == nil {
		t.Fatalf("expected error for insufficient stock")
	}
	if len(store.orders) != 0 || len(store.orderItems) != 0 {
		t.Fatalf("expected no persisted order or items")
	}
	if store.cartCleared {
		t.Fatalf("expected cart not cleared")
	}
}

type staleAtomicProductRepo struct {
	product      *domain.Product
	atomicCalled bool
}

func (r *staleAtomicProductRepo) Create(product *domain.Product) error { return nil }
func (r *staleAtomicProductRepo) Update(product *domain.Product) error { return nil }
func (r *staleAtomicProductRepo) Delete(id int64) error                { return nil }
func (r *staleAtomicProductRepo) FindByID(id int64) (*domain.Product, error) {
	if r.product == nil {
		return nil, errors.New("product not found")
	}
	return r.product, nil
}
func (r *staleAtomicProductRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (r *staleAtomicProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (r *staleAtomicProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
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
func (r *staleAtomicProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, nil
}
func (r *staleAtomicProductRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *staleAtomicProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (r *staleAtomicProductRepo) SearchCount(query string) (int64, error)  { return 0, nil }
func (r *staleAtomicProductRepo) UpdateStock(id int64, quantity int) error { return nil }
func (r *staleAtomicProductRepo) IncrementViews(id int64) error            { return nil }
func (r *staleAtomicProductRepo) IncrementSales(id int64, quantity int) error {
	return nil
}
func (r *staleAtomicProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	r.atomicCalled = true
	return nil
}

func TestOrderServiceCreateOrderUsesAtomicStockUpdateWhenStaleRead(t *testing.T) {
	orderRepo := &fakeOrderCreateRepo{}
	productRepo := &staleAtomicProductRepo{product: &domain.Product{
		ID:            50,
		Name:          "Cat Collar",
		SKU:           "COLLAR-001",
		StockQuantity: 0,
	}}
	inventoryRepo := &fakeInventoryRepo{}
	cartRepo := &fakeCartRepo{cart: &domain.Cart{
		ID:     500,
		UserID: 5,
		Items: []domain.CartItem{{
			ProductID: 50,
			Quantity:  1,
			Price:     8,
		}},
	}}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(5, "", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if !productRepo.atomicCalled {
		t.Fatalf("expected atomic stock update to be used")
	}
	if len(inventoryRepo.logs) != 1 {
		t.Fatalf("expected inventory log to be created")
	}
}

type orderTestStore struct {
	orders         []*domain.Order
	orderItems     []*domain.OrderItem
	inventoryLogs  []*domain.InventoryLog
	cartCleared    bool
	cart           *domain.Cart
	products       map[int64]*domain.Product
	updatedStock   map[int64]int
	updateStockErr error
	cartClearErr   error
	sales          map[int64]int
}

func (s *orderTestStore) clone() *orderTestStore {
	clone := &orderTestStore{
		orders:         append([]*domain.Order{}, s.orders...),
		orderItems:     append([]*domain.OrderItem{}, s.orderItems...),
		inventoryLogs:  append([]*domain.InventoryLog{}, s.inventoryLogs...),
		cartCleared:    s.cartCleared,
		updateStockErr: s.updateStockErr,
		cartClearErr:   s.cartClearErr,
		updatedStock:   map[int64]int{},
		products:       map[int64]*domain.Product{},
		sales:          map[int64]int{},
	}
	for id, product := range s.products {
		copyProduct := *product
		clone.products[id] = &copyProduct
	}
	for id, qty := range s.updatedStock {
		clone.updatedStock[id] = qty
	}
	for id, qty := range s.sales {
		clone.sales[id] = qty
	}
	if s.cart != nil {
		cartCopy := *s.cart
		cartCopy.Items = append([]domain.CartItem{}, s.cart.Items...)
		clone.cart = &cartCopy
	}
	return clone
}

type txOrderRepo struct {
	store             *orderTestStore
	transactionCalled bool
}

func (r *txOrderRepo) Create(order *domain.Order) error {
	r.store.orders = append(r.store.orders, order)
	return nil
}
func (r *txOrderRepo) Update(order *domain.Order) error                    { return nil }
func (r *txOrderRepo) FindByID(id int64) (*domain.Order, error)            { return nil, nil }
func (r *txOrderRepo) FindByOrderNo(orderNo string) (*domain.Order, error) { return nil, nil }
func (r *txOrderRepo) FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	return nil, nil
}
func (r *txOrderRepo) CountByUserID(userID int64) (int64, error) { return 0, nil }
func (r *txOrderRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error) {
	return nil, nil
}
func (r *txOrderRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *txOrderRepo) UpdateStatus(id int64, status string) error          { return nil }
func (r *txOrderRepo) CreateOrderItem(item *domain.OrderItem) error {
	r.store.orderItems = append(r.store.orderItems, item)
	return nil
}
func (r *txOrderRepo) Transaction(fn func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, idempotencyRepo repository.OrderIdempotencyRepository, paymentRepo repository.PaymentRepository) error) error {
	r.transactionCalled = true
	clone := r.store.clone()
	orderRepo := &txOrderRepo{store: clone}
	cartRepo := &txCartRepo{store: clone}
	productRepo := &txProductRepo{store: clone}
	inventoryRepo := &txInventoryRepo{store: clone}
	if err := fn(orderRepo, cartRepo, productRepo, inventoryRepo, nil, nil); err != nil {
		return err
	}
	*r.store = *clone
	return nil
}

type txCartRepo struct {
	store *orderTestStore
}

func (r *txCartRepo) FindByUserID(userID int64) (*domain.Cart, error) { return r.store.cart, nil }
func (r *txCartRepo) Create(cart *domain.Cart) error                  { return nil }
func (r *txCartRepo) Update(cart *domain.Cart) error                  { return nil }
func (r *txCartRepo) Delete(id int64) error                           { return nil }
func (r *txCartRepo) AddItem(item *domain.CartItem) error             { return nil }
func (r *txCartRepo) UpdateItem(item *domain.CartItem) error          { return nil }
func (r *txCartRepo) RemoveItem(cartID, productID int64) error        { return nil }
func (r *txCartRepo) ClearCart(cartID int64) error {
	if r.store.cartClearErr != nil {
		return r.store.cartClearErr
	}
	r.store.cartCleared = true
	return nil
}

type txProductRepo struct {
	store *orderTestStore
}

func (r *txProductRepo) Create(product *domain.Product) error { return nil }
func (r *txProductRepo) Update(product *domain.Product) error { return nil }
func (r *txProductRepo) Delete(id int64) error                { return nil }
func (r *txProductRepo) FindByID(id int64) (*domain.Product, error) {
	product, ok := r.store.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}
func (r *txProductRepo) FindBySKU(sku string) (*domain.Product, error) { return nil, nil }
func (r *txProductRepo) FindBySlug(slug string) (*domain.Product, error) {
	return nil, nil
}
func (r *txProductRepo) GetByIDs(ids []int64) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0, len(ids))
	for _, id := range ids {
		product, ok := r.store.products[id]
		if !ok {
			return nil, errors.New("product not found")
		}
		products = append(products, product)
	}
	return products, nil
}
func (r *txProductRepo) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	return nil, nil
}
func (r *txProductRepo) Count(filters map[string]interface{}) (int64, error) { return 0, nil }
func (r *txProductRepo) Search(query string, offset, limit int) ([]*domain.Product, error) {
	return nil, nil
}
func (r *txProductRepo) SearchCount(query string) (int64, error) { return 0, nil }
func (r *txProductRepo) UpdateStock(id int64, quantity int) error {
	if r.store.updateStockErr != nil {
		return r.store.updateStockErr
	}
	if r.store.updatedStock == nil {
		r.store.updatedStock = map[int64]int{}
	}
	r.store.updatedStock[id] = quantity
	return nil
}
func (r *txProductRepo) UpdateStockWithDelta(id int64, delta int) error {
	product, ok := r.store.products[id]
	if !ok {
		return errors.New("product not found")
	}
	newStock := product.StockQuantity + delta
	if newStock < 0 {
		return errors.New("insufficient stock")
	}
	product.StockQuantity = newStock
	return r.UpdateStock(id, newStock)
}
func (r *txProductRepo) IncrementViews(id int64) error { return nil }
func (r *txProductRepo) IncrementSales(id int64, quantity int) error {
	if r.store.sales == nil {
		r.store.sales = map[int64]int{}
	}
	r.store.sales[id] += quantity
	return nil
}

type txInventoryRepo struct {
	store *orderTestStore
}

func (r *txInventoryRepo) Create(log *domain.InventoryLog) error {
	r.store.inventoryLogs = append(r.store.inventoryLogs, log)
	return nil
}
func (r *txInventoryRepo) FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error) {
	return nil, nil
}
func (r *txInventoryRepo) CountByProductID(productID int64) (int64, error) { return 0, nil }

func TestOrderServiceCreateOrderRollsBackOnInventoryFailure(t *testing.T) {
	store := &orderTestStore{
		products: map[int64]*domain.Product{
			30: {
				ID:            30,
				Name:          "Cat Bed",
				SKU:           "BED-001",
				StockQuantity: 5,
			},
		},
		cart: &domain.Cart{
			ID:     300,
			UserID: 3,
			Items: []domain.CartItem{{
				ProductID: 30,
				Quantity:  1,
				Price:     99,
			}},
		},
		updateStockErr: errors.New("update stock failed"),
	}
	orderRepo := &txOrderRepo{store: store}
	productRepo := &txProductRepo{store: store}
	inventoryRepo := &txInventoryRepo{store: store}
	cartRepo := &txCartRepo{store: store}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(3, "", "ship", "bill", "", "card")
	if err == nil {
		t.Fatalf("expected error for inventory update failure")
	}
	if !orderRepo.transactionCalled {
		t.Fatalf("expected transaction to be used")
	}
	if len(store.orders) != 0 || len(store.orderItems) != 0 {
		t.Fatalf("expected no persisted order or items after rollback")
	}
	if store.cartCleared {
		t.Fatalf("expected cart not cleared after rollback")
	}
}

func TestOrderServiceCreateOrderRollsBackOnCartClearFailure(t *testing.T) {
	store := &orderTestStore{
		products: map[int64]*domain.Product{
			40: {
				ID:            40,
				Name:          "Cat Tower",
				SKU:           "TOWER-001",
				StockQuantity: 3,
			},
		},
		cart: &domain.Cart{
			ID:     400,
			UserID: 4,
			Items: []domain.CartItem{{
				ProductID: 40,
				Quantity:  1,
				Price:     120,
			}},
		},
		cartClearErr: errors.New("clear cart failed"),
	}
	orderRepo := &txOrderRepo{store: store}
	productRepo := &txProductRepo{store: store}
	inventoryRepo := &txInventoryRepo{store: store}
	cartRepo := &txCartRepo{store: store}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(4, "", "ship", "bill", "", "card")
	if err == nil {
		t.Fatalf("expected error for cart clear failure")
	}
	if !orderRepo.transactionCalled {
		t.Fatalf("expected transaction to be used")
	}
	if len(store.orders) != 0 || len(store.orderItems) != 0 {
		t.Fatalf("expected no persisted order or items after rollback")
	}
	if len(store.inventoryLogs) != 0 {
		t.Fatalf("expected no inventory logs after rollback")
	}
	if store.cartCleared {
		t.Fatalf("expected cart not cleared after rollback")
	}
}

func TestOrderServiceCreateOrderIncrementsProductSales(t *testing.T) {
	store := &orderTestStore{
		products: map[int64]*domain.Product{
			60: {
				ID:            60,
				Name:          "Cat Harness",
				SKU:           "HARNESS-001",
				StockQuantity: 10,
			},
			61: {
				ID:            61,
				Name:          "Cat Leash",
				SKU:           "LEASH-001",
				StockQuantity: 8,
			},
		},
		cart: &domain.Cart{
			ID:     600,
			UserID: 6,
			Items: []domain.CartItem{{
				ProductID: 60,
				Quantity:  2,
				Price:     20,
			}, {
				ProductID: 61,
				Quantity:  1,
				Price:     15,
			}},
		},
	}
	orderRepo := &txOrderRepo{store: store}
	productRepo := &txProductRepo{store: store}
	inventoryRepo := &txInventoryRepo{store: store}
	cartRepo := &txCartRepo{store: store}

	svc := NewOrderService(orderRepo, cartRepo, productRepo, inventoryRepo, nil)
	_, err := svc.CreateOrder(6, "", "ship", "bill", "", "card")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if store.sales[60] != 2 {
		t.Fatalf("expected product 60 sales to be incremented")
	}
	if store.sales[61] != 1 {
		t.Fatalf("expected product 61 sales to be incremented")
	}
}
