package repository

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int64) error
	FindByID(id int64) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	List(offset, limit int) ([]*domain.User, error)
	Count() (int64, error)
}

type ProductRepository interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id int64) error
	FindByID(id int64) (*domain.Product, error)
	FindBySKU(sku string) (*domain.Product, error)
	FindBySlug(slug string) (*domain.Product, error)
	GetByIDs(ids []int64) ([]*domain.Product, error)
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error)
	Count(filters map[string]interface{}) (int64, error)
	Search(query string, offset, limit int) ([]*domain.Product, error)
	SearchCount(query string) (int64, error)
	UpdateStock(id int64, quantity int) error
	UpdateStockWithDelta(id int64, delta int) error
	IncrementViews(id int64) error
	IncrementSales(id int64, quantity int) error
}

type CategoryRepository interface {
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(id int64) error
	FindByID(id int64) (*domain.Category, error)
	FindBySlug(slug string) (*domain.Category, error)
	List(offset, limit int) ([]*domain.Category, error)
	Count() (int64, error)
	Tree() ([]*domain.Category, error)
}

type CartRepository interface {
	FindByUserID(userID int64) (*domain.Cart, error)
	Create(cart *domain.Cart) error
	Update(cart *domain.Cart) error
	Delete(id int64) error
	AddItem(item *domain.CartItem) error
	UpdateItem(item *domain.CartItem) error
	RemoveItem(cartID, productID int64) error
	ClearCart(cartID int64) error
}

type OrderRepository interface {
	Create(order *domain.Order) error
	Update(order *domain.Order) error
	FindByID(id int64) (*domain.Order, error)
	FindByOrderNo(orderNo string) (*domain.Order, error)
	FindByUserID(userID int64, offset, limit int) ([]*domain.Order, error)
	CountByUserID(userID int64) (int64, error)
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Order, error)
	Count(filters map[string]interface{}) (int64, error)
	UpdateStatus(id int64, status string) error
	CreateOrderItem(item *domain.OrderItem) error
}

type ShipmentRepository interface {
	Create(shipment *domain.Shipment) error
	FindByOrderID(orderID int64) (*domain.Shipment, error)
	Update(shipment *domain.Shipment) error
}

type OrderStatusLogRepository interface {
	Create(log *domain.OrderStatusLog) error
}

type OrderIdempotencyRepository interface {
	Create(record *domain.OrderIdempotency) error
	FindByUserIDAndIdempotencyKey(userID int64, key string) (*domain.OrderIdempotency, error)
	Update(record *domain.OrderIdempotency) error
}

type IdempotencyKeyRepository interface {
	Create(record *domain.IdempotencyKey) error
	FindByUserIDAndKey(userID int64, key string) (*domain.IdempotencyKey, error)
	Update(record *domain.IdempotencyKey) error
}

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	Update(payment *domain.Payment) error
	FindByID(id int64) (*domain.Payment, error)
	FindByOrderID(orderID int64) ([]*domain.Payment, error)
	FindByTransactionID(transactionID string) (*domain.Payment, error)
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, error)
	Count(filters map[string]interface{}) (int64, error)
}

type PaymentRefundRepository interface {
	Create(refund *domain.PaymentRefund) error
}

type PaymentEventRepository interface {
	Create(event *domain.PaymentEvent) error
}

type InventoryRepository interface {
	Create(log *domain.InventoryLog) error
	FindByProductID(productID int64, offset, limit int) ([]*domain.InventoryLog, error)
	CountByProductID(productID int64) (int64, error)
}

type CheckoutSessionRepository interface {
	Create(session *domain.CheckoutSession) error
	Update(session *domain.CheckoutSession) error
	FindByID(id string) (*domain.CheckoutSession, error)
	Delete(id string) error
}

type PaymentHandlerRepository interface {
	Create(handler *domain.PaymentHandler) error
	Update(handler *domain.PaymentHandler) error
	FindByID(id int64) (*domain.PaymentHandler, error)
	FindByName(name string) (*domain.PaymentHandler, error)
	List() ([]*domain.PaymentHandler, error)
}

type OAuthClientRepository interface {
	Create(client *domain.OAuthClient) error
	List(offset, limit int) ([]*domain.OAuthClient, error)
	Count() (int64, error)
}

type OAuthTokenRepository interface {
	Create(token *domain.OAuthToken) error
	FindByToken(token string) (*domain.OAuthToken, error)
	Revoke(token string, revokedAt time.Time) error
}

type UCPWebhookEventRepository interface {
	Create(event *domain.UCPWebhookEvent) error
	FindByEventID(eventID string) (*domain.UCPWebhookEvent, error)
	UpdateStatus(eventID string, status string) error
	MarkProcessed(eventID string) error
}

type UCPWebhookAuditRepository interface {
	Create(audit *domain.UCPWebhookAudit) error
	List(offset, limit int) ([]*domain.UCPWebhookAudit, error)
	Count() (int64, error)
}

type UCPWebhookReplayRepository interface {
	FindByHash(hash string) (*domain.UCPWebhookReplay, error)
	Create(replay *domain.UCPWebhookReplay) error
}

type UCPWebhookQueueRepository interface {
	Create(job *domain.UCPWebhookJob) error
	ListDue(limit int) ([]*domain.UCPWebhookJob, error)
	Update(job *domain.UCPWebhookJob) error
	List(offset, limit int) ([]*domain.UCPWebhookJob, error)
	Count() (int64, error)
	FindByID(id int64) (*domain.UCPWebhookJob, error)
}

type UCPWebhookAlertRepository interface {
	Create(alert *domain.UCPWebhookAlert) error
	List(offset, limit int) ([]*domain.UCPWebhookAlert, error)
	Count() (int64, error)
	ExistsRecent(eventID, reason string, window time.Duration) (bool, error)
}

type WebhookDLQRepository interface {
	Create(item *domain.WebhookDLQ) error
	FindByID(id int64) (*domain.WebhookDLQ, error)
	List(offset, limit int) ([]*domain.WebhookDLQ, error)
	Count() (int64, error)
}

type WebhookReplayLogRepository interface {
	Create(item *domain.WebhookReplayLog) error
}

type Repositories struct {
	User             UserRepository
	Product          ProductRepository
	Category         CategoryRepository
	Cart             CartRepository
	Order            OrderRepository
	OrderIdempotency OrderIdempotencyRepository
	IdempotencyKey   IdempotencyKeyRepository
	Shipment         ShipmentRepository
	OrderStatusLog   OrderStatusLogRepository
	Payment          PaymentRepository
	PaymentRefund    PaymentRefundRepository
	PaymentEvent     PaymentEventRepository
	Inventory        InventoryRepository
	Checkout         CheckoutSessionRepository
	Handler          PaymentHandlerRepository
	OAuthClient      OAuthClientRepository
	OAuthToken       OAuthTokenRepository
	Webhook          UCPWebhookEventRepository
	WebhookAudit     UCPWebhookAuditRepository
	WebhookReplay    UCPWebhookReplayRepository
	WebhookQueue     UCPWebhookQueueRepository
	WebhookAlert     UCPWebhookAlertRepository
	WebhookDLQ       WebhookDLQRepository
	WebhookReplayLog WebhookReplayLogRepository
}

func NewRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		Product:          NewProductRepository(db),
		Category:         NewCategoryRepository(db),
		Cart:             NewCartRepository(db),
		Order:            NewOrderRepository(db),
		OrderIdempotency: NewOrderIdempotencyRepository(db),
		IdempotencyKey:   NewIdempotencyKeyRepository(db),
		Shipment:         NewShipmentRepository(db),
		OrderStatusLog:   NewOrderStatusLogRepository(db),
		Payment:          NewPaymentRepository(db),
		PaymentRefund:    NewPaymentRefundRepository(db),
		PaymentEvent:     NewPaymentEventRepository(db),
		Inventory:        NewInventoryRepository(db),
		Checkout:         NewCheckoutSessionRepository(db),
		Handler:          NewPaymentHandlerRepository(db),
		OAuthClient:      NewOAuthClientRepository(db),
		OAuthToken:       NewOAuthTokenRepository(db),
		Webhook:          NewUCPWebhookEventRepository(db),
		WebhookAudit:     NewUCPWebhookAuditRepository(db),
		WebhookReplay:    NewUCPWebhookReplayRepository(db),
		WebhookQueue:     NewUCPWebhookQueueRepository(db),
		WebhookAlert:     NewUCPWebhookAlertRepository(db),
		WebhookDLQ:       NewWebhookDLQRepository(db),
		WebhookReplayLog: NewWebhookReplayLogRepository(db),
	}
}
