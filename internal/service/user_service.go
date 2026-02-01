package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

type Services struct {
	User            *UserService
	Product         *ProductService
	Category        *CategoryService
	Cart            *CartService
	Order           *OrderService
	Payment         *PaymentService
	Inventory       *InventoryService
	Checkout        *CheckoutSessionService
	Handler         *PaymentHandlerService
	Webhook         *WebhookEventService
	UCPOrder        *UCPOrderService
	WebhookAudit    *WebhookAuditService
	WebhookReplay   *WebhookReplayService
	WebhookQueue    *WebhookQueueService
	WebhookDLQ      *WebhookDLQService
	WebhookAlert    *WebhookAlertService
	OAuthClient     *OAuthClientService
	OAuthToken      *OAuthTokenService
	OAuthClientRepo repository.OAuthClientRepository
	OAuthTokenRepo  repository.OAuthTokenRepository
	Promotion       *PromotionService
}

func NewServices(repos *repository.Repositories, redis *redis.Client) *Services {
	orderService := NewOrderService(repos.Order, repos.Cart, repos.Product, repos.Inventory, repos.OrderIdempotency)
	webhookQueue := NewWebhookQueueService(repos.WebhookQueue)
	orderService.SetWebhookQueue(webhookQueue)
	orderService.SetShipmentRepo(repos.Shipment)
	orderService.SetStatusLogRepo(repos.OrderStatusLog)
	paymentService := NewPaymentServiceWithDeps(repos.Payment, repos.Order, repos.PaymentRefund, repos.PaymentEvent)
	webhookDLQ := NewWebhookDLQService(webhookQueue, repos.WebhookDLQ)
	oauthClient := NewOAuthClientService(repos.OAuthClient)
	oauthToken := NewOAuthTokenService(repos.OAuthClient, repos.OAuthToken)
	taxShipping := NewTaxShippingService(repos.TaxRule, repos.ShippingRule)
	checkoutService := NewCheckoutSessionService(repos.Checkout)
	checkoutService.SetTaxShippingService(taxShipping)
	promotionService := NewPromotionService(repos.Coupon)

	return &Services{
		User:            NewUserService(repos.User),
		Product:         NewProductService(repos.Product, repos.Inventory, redis),
		Category:        NewCategoryService(repos.Category),
		Cart:            NewCartService(repos.Cart, repos.Product),
		Order:           orderService,
		Payment:         paymentService,
		Inventory:       NewInventoryService(repos.Product, repos.Inventory),
		Checkout:        checkoutService,
		Promotion:       promotionService,
		Handler:         NewPaymentHandlerService(repos.Handler),
		Webhook:         NewWebhookEventService(repos.Webhook),
		UCPOrder:        NewUCPOrderService(repos.Order, repos.Payment),
		WebhookAudit:    NewWebhookAuditService(repos.WebhookAudit),
		WebhookReplay:   NewWebhookReplayService(repos.WebhookReplay),
		WebhookQueue:    webhookQueue,
		WebhookDLQ:      webhookDLQ,
		WebhookAlert:    NewWebhookAlertService(repos.WebhookAlert, repos.Webhook),
		OAuthClient:     oauthClient,
		OAuthToken:      oauthToken,
		OAuthClientRepo: repos.OAuthClient,
		OAuthTokenRepo:  repos.OAuthToken,
	}
}

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Create(user)
}

func (s *UserService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id int64) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) ListUsers(offset, limit int) ([]*domain.User, int64, error) {
	users, err := s.userRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (s *UserService) DeleteUser(id int64) error {
	return s.userRepo.Delete(id)
}
