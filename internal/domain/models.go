package domain

import "time"

type User struct {
	ID           int64  `gorm:"primary_key"`
	Username     string `gorm:"unique_index;not null"`
	Email        string `gorm:"unique_index;not null"`
	PasswordHash string `gorm:"not null"`
	Phone        string
	Avatar       string
	Status       int    `gorm:"default:1;check:status IN (0, 1)"`
	Role         string `gorm:"default:'user';check:role IN ('user', 'admin', 'super_admin')"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Category struct {
	ID          int64  `gorm:"primary_key"`
	Name        string `gorm:"not null"`
	Slug        string `gorm:"unique_index;not null"`
	Description string
	ParentID    *int64 `gorm:"index"`
	SortOrder   int    `gorm:"default:0"`
	Status      int    `gorm:"default:1;check:status IN (0, 1)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	ID                int64  `gorm:"primary_key"`
	Name              string `gorm:"not null"`
	Slug              string `gorm:"unique_index;not null"`
	Description       string
	Price             float64 `gorm:"type:decimal(10,2);not null"`
	ComparePrice      float64 `gorm:"type:decimal(10,2)"`
	SKU               string  `gorm:"unique_index;not null"`
	StockQuantity     int     `gorm:"default:0"`
	LowStockThreshold int     `gorm:"default:10"`
	CategoryID        *int64  `gorm:"index"`
	Images            string
	Status            int     `gorm:"default:1;check:status IN (0, 1, 2)"`
	Featured          bool    `gorm:"default:false"`
	Views             int     `gorm:"default:0"`
	Sales             int     `gorm:"default:0"`
	Weight            float64 `gorm:"type:decimal(10,2)"`
	Dimensions        string
	MetaTitle         string
	MetaDescription   string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Cart struct {
	ID        int64 `gorm:"primary_key"`
	UserID    int64 `gorm:"unique_index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User       `gorm:"foreignkey:UserID"`
	Items     []CartItem `gorm:"foreignkey:CartID"`
}

type CartItem struct {
	ID        int64   `gorm:"primary_key"`
	CartID    int64   `gorm:"index;not null"`
	ProductID int64   `gorm:"index;not null"`
	Quantity  int     `gorm:"not null;check:quantity > 0"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID              int64   `gorm:"primary_key"`
	OrderNo         string  `gorm:"unique_index;not null"`
	UserID          *int64  `gorm:"index"`
	Status          string  `gorm:"default:'pending';check:status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled', 'refunded')"`
	Subtotal        float64 `gorm:"type:decimal(10,2);not null"`
	ShippingFee     float64 `gorm:"type:decimal(10,2);default:0"`
	Tax             float64 `gorm:"type:decimal(10,2);default:0"`
	Discount        float64 `gorm:"type:decimal(10,2);default:0"`
	Total           float64 `gorm:"type:decimal(10,2);not null"`
	Currency        string  `gorm:"default:'CNY'"`
	PaymentMethod   string
	PaymentStatus   string `gorm:"default:'unpaid'"`
	PaymentTime     *time.Time
	ShippingAddress string
	BillingAddress  string
	Notes           string
	ShippedAt       *time.Time
	DeliveredAt     *time.Time
	CancelledAt     *time.Time
	RefundedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	User            User        `gorm:"foreignkey:UserID"`
	Items           []OrderItem `gorm:"foreignkey:OrderID"`
	Payments        []Payment   `gorm:"foreignkey:OrderID"`
}

type OrderIdempotency struct {
	ID             int64  `gorm:"primary_key"`
	UserID         int64  `gorm:"index;not null"`
	IdempotencyKey string `gorm:"not null"`
	OrderID        *int64
	Status         string `gorm:"not null;default:'pending'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type IdempotencyKey struct {
	ID               int64 `gorm:"primary_key"`
	UserID           int64 `gorm:"not null"`
	Key              string
	RequestHash      string
	ResponseSnapshot *string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OrderItem struct {
	ID          int64  `gorm:"primary_key"`
	OrderID     int64  `gorm:"index;not null"`
	ProductID   *int64 `gorm:"index"`
	ProductName string `gorm:"not null"`
	SKU         string
	Quantity    int     `gorm:"not null;check:quantity > 0"`
	UnitPrice   float64 `gorm:"type:decimal(10,2);not null"`
	TotalPrice  float64 `gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time
}

type Payment struct {
	ID             int64 `gorm:"primary_key"`
	OrderID        int64 `gorm:"index;not null"`
	UserID         *int64
	Amount         float64 `gorm:"type:decimal(10,2);not null"`
	PaymentMethod  string  `gorm:"not null"`
	TransactionID  string
	Status         string `gorm:"default:'pending'"`
	ErrorMessage   string
	PaymentPayload string `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PaymentRefund struct {
	ID          int64 `gorm:"primary_key"`
	PaymentID   int64 `gorm:"not null"`
	Amount      float64
	Status      string
	Reason      string
	ExternalRef *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PaymentEvent struct {
	ID        int64 `gorm:"primary_key"`
	PaymentID int64 `gorm:"not null"`
	EventType string
	Payload   *string
	CreatedAt time.Time
}

type OrderStatusLog struct {
	ID         int64 `gorm:"primary_key"`
	OrderID    int64 `gorm:"not null"`
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  time.Time
}

type Shipment struct {
	ID          int64 `gorm:"primary_key"`
	OrderID     int64 `gorm:"not null"`
	Carrier     string
	TrackingNo  string
	Status      string
	ShippedAt   *time.Time
	DeliveredAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryLog struct {
	ID             int64  `gorm:"primary_key"`
	ProductID      int64  `gorm:"index;not null"`
	QuantityChange int    `gorm:"not null"`
	Type           string `gorm:"not null;check:type IN ('in', 'out', 'adjust')"`
	ReferenceID    string
	ReferenceType  string
	Notes          string
	CreatedAt      time.Time
}

type CheckoutSession struct {
	ID          string `gorm:"primary_key"`
	Status      string `gorm:"not null"`
	Currency    string `gorm:"not null"`
	LineItems   string `gorm:"type:jsonb;not null"`
	Totals      string `gorm:"type:jsonb;not null"`
	Buyer       string `gorm:"type:jsonb"`
	Messages    string `gorm:"type:jsonb"`
	Links       string `gorm:"type:jsonb"`
	ContinueURL string `gorm:"type:text"`
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PaymentHandler struct {
	ID           int64  `gorm:"primary_key"`
	Name         string `gorm:"not null"`
	Version      string `gorm:"not null"`
	Spec         string `gorm:"type:text;not null"`
	ConfigSchema string `gorm:"type:text;not null"`
	Config       string `gorm:"type:jsonb;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UCPWebhookEvent struct {
	ID          int64  `gorm:"primary_key"`
	EventID     string `gorm:"unique_index;not null"`
	EventType   string `gorm:"not null"`
	OrderID     string `gorm:"not null"`
	Status      string `gorm:"not null"`
	PayloadHash string `gorm:"not null"`
	ReceivedAt  time.Time
	ProcessedAt *time.Time
}

type UCPWebhookAudit struct {
	ID              int64  `gorm:"primary_key"`
	EventID         string `gorm:"index"`
	Reason          string `gorm:"not null"`
	SignatureHeader string
	KeyID           string
	PayloadHash     string `gorm:"not null"`
	CreatedAt       time.Time
}

type UCPWebhookReplay struct {
	ID          int64     `gorm:"primary_key"`
	PayloadHash string    `gorm:"unique_index;not null"`
	ExpiresAt   time.Time `gorm:"index"`
	CreatedAt   time.Time
}

type UCPWebhookJob struct {
	ID            int64     `gorm:"primary_key"`
	EventID       string    `gorm:"index;not null"`
	Payload       string    `gorm:"type:text;not null"`
	Status        string    `gorm:"not null"`
	Attempts      int       `gorm:"not null"`
	NextRetryAt   time.Time `gorm:"index"`
	LastError     string
	LastAttemptAt time.Time
	CreatedAt     time.Time
}

type UCPWebhookAlert struct {
	ID        int64  `gorm:"primary_key"`
	EventID   string `gorm:"index"`
	Reason    string `gorm:"not null"`
	Details   string
	Attempts  int
	CreatedAt time.Time
}

type WebhookDLQ struct {
	ID        int64 `gorm:"primary_key"`
	JobID     int64
	Reason    string
	Payload   string
	CreatedAt time.Time
}

type WebhookReplayLog struct {
	ID       int64 `gorm:"primary_key"`
	JobID    int64
	ReplayAt time.Time
	Result   string
}

type OAuthClient struct {
	ID         int64 `gorm:"primary_key"`
	ClientID   string
	SecretHash string
	Scopes     string
	Status     string
	CreatedAt  time.Time
}

type OAuthToken struct {
	ID        int64 `gorm:"primary_key"`
	Token     string
	ClientID  string
	UserID    *int64
	Scopes    string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type TaxRule struct {
	ID          int64 `gorm:"primary_key"`
	Region      string
	Category    string
	Rate        float64
	EffectiveAt time.Time
}

type ShippingRule struct {
	ID            int64 `gorm:"primary_key"`
	Region        string
	Method        string
	BaseAmount    float64
	PerItemAmount float64
}
