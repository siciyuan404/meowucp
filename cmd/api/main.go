package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/api"
	"github.com/meowucp/internal/middleware"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/service"
	ucpapi "github.com/meowucp/internal/ucp/api"
	ucpmodel "github.com/meowucp/internal/ucp/model"
	"github.com/meowucp/internal/ucp/security"
	"github.com/meowucp/pkg/config"
	"github.com/meowucp/pkg/database"
	"github.com/meowucp/pkg/redis"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient, err := redis.NewClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PoolSize,
	)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, redisClient)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORS())

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)
	ucpProfileHandler := ucpapi.NewProfileHandler(services)
	ucpCheckoutHandler := ucpapi.NewCheckoutHandlerWithConfig(services, ucpapi.CheckoutHandlerConfig{
		Links:           buildUCPLinks(cfg.UCP.Links),
		ContinueURLBase: cfg.UCP.ContinueURLBase,
	})
	ucpVerifier := security.NewJWKVerifier(cfg.UCP.Webhook.JWKSetURL, cfg.UCP.Webhook.ClockSkewSeconds)
	ucpVerifier.SetSkipVerify(cfg.UCP.Webhook.SkipSignatureVerify)
	ucpVerifier.SetNonceStore(ucpapi.NewWebhookReplayNonceStore(services.WebhookReplay))
	ucpOrderWebhookHandler := ucpapi.NewOrderWebhookHandlerWithVerifier(services, ucpVerifier)
	paymentCallbackHandler := api.NewPaymentCallbackHandler(services.Payment, services.Order)
	paymentRefundHandler := api.NewPaymentRefundHandler(services.Payment)
	oauthMetadataHandler := api.NewOAuthMetadataHandler()
	oauthTokenHandler := api.NewOAuthTokenHandler()
	oauthAuthorizeHandler := api.NewOAuthAuthorizeHandler()
	adminOrderWebhookHandler := api.NewAdminOrderWebhookHandler(services.Order, services.WebhookQueue, api.AdminOrderWebhookConfig{
		DeliveryURL: cfg.UCP.Webhook.DeliveryURL,
		Timeout:     time.Duration(cfg.UCP.Webhook.DeliveryTimeoutSec) * time.Second,
	})
	adminWebhookDLQHandler := api.NewAdminWebhookDLQHandler(services.WebhookDLQ)
	adminOAuthClientHandler := api.NewAdminOAuthClientHandler(services.OAuthClient)

	apiGroup := r.Group("/api/v1")
	{
		public := apiGroup.Group("/public")
		{
			public.POST("/register", func(c *gin.Context) {
				userHandler := api.NewUserHandler(services.User)
				userHandler.Register(c)
			})
			public.POST("/login", func(c *gin.Context) {
				userHandler := api.NewUserHandler(services.User)
				userHandler.Login(c)
			})
		}

		user := apiGroup.Group("/user")
		user.Use(authMiddleware.Auth())
		{
			user.GET("/me", func(c *gin.Context) {
				userHandler := api.NewUserHandler(services.User)
				userHandler.GetCurrentUser(c)
			})
			user.PUT("/me", func(c *gin.Context) {
				userHandler := api.NewUserHandler(services.User)
				userHandler.UpdateCurrentUser(c)
			})
		}

		orderHandler := api.NewOrderHandler(services.Order)
		apiGroup.POST("/orders", func(c *gin.Context) {
			orderHandler.Create(c)
		})
		apiGroup.POST("/payment/callback", func(c *gin.Context) {
			paymentCallbackHandler.Handle(c)
		})
		apiGroup.POST("/payments/:id/refund", func(c *gin.Context) {
			paymentRefundHandler.Create(c)
		})

		admin := apiGroup.Group("/admin")
		admin.Use(authMiddleware.Auth(), authMiddleware.AdminOnly())
		{
			admin.GET("/users", func(c *gin.Context) {
				userHandler := api.NewUserHandler(services.User)
				userHandler.ListUsers(c)
			})
			admin.GET("/ucp/webhook-audits", func(c *gin.Context) {
				webhookAuditHandler := api.NewWebhookAuditHandler(services.WebhookAudit)
				webhookAuditHandler.List(c)
			})
			admin.GET("/ucp/webhook-alerts", func(c *gin.Context) {
				webhookAlertHandler := api.NewWebhookAlertHandler(services.WebhookAlert)
				webhookAlertHandler.List(c)
			})
			admin.GET("/ucp/webhook-jobs", func(c *gin.Context) {
				webhookJobHandler := api.NewWebhookJobHandler(services.WebhookQueue)
				webhookJobHandler.List(c)
			})
			admin.POST("/ucp/webhook-jobs/:id/retry", func(c *gin.Context) {
				webhookJobHandler := api.NewWebhookJobHandler(services.WebhookQueue)
				webhookJobHandler.Retry(c)
			})

			adminProductHandler := api.NewAdminProductHandler(services.Product)
			admin.GET("/products", func(c *gin.Context) {
				adminProductHandler.List(c)
			})
			admin.POST("/products", func(c *gin.Context) {
				adminProductHandler.Create(c)
			})
			admin.GET("/products/:id", func(c *gin.Context) {
				adminProductHandler.Get(c)
			})
			admin.PUT("/products/:id", func(c *gin.Context) {
				adminProductHandler.Update(c)
			})
			admin.PATCH("/products/:id/status", func(c *gin.Context) {
				adminProductHandler.UpdateStatus(c)
			})

			adminCategoryHandler := api.NewAdminCategoryHandler(services.Category)
			admin.GET("/categories", func(c *gin.Context) {
				adminCategoryHandler.List(c)
			})
			admin.POST("/categories", func(c *gin.Context) {
				adminCategoryHandler.Create(c)
			})
			admin.PUT("/categories/:id", func(c *gin.Context) {
				adminCategoryHandler.Update(c)
			})

			adminInventoryHandler := api.NewAdminInventoryHandler(adminInventoryServiceAdapter{svc: services.Inventory})
			admin.POST("/inventory/adjust", func(c *gin.Context) {
				adminInventoryHandler.Adjust(c)
			})
			admin.GET("/inventory/logs", func(c *gin.Context) {
				adminInventoryHandler.Logs(c)
			})

			adminOrderHandler := api.NewAdminOrderHandler(services.Order)
			admin.GET("/orders", func(c *gin.Context) {
				adminOrderHandler.List(c)
			})
			admin.GET("/orders/:id", func(c *gin.Context) {
				adminOrderHandler.Get(c)
			})
			admin.POST("/orders/:id/webhook", func(c *gin.Context) {
				adminOrderWebhookHandler.Trigger(c)
			})
			admin.POST("/orders/:id/ship", func(c *gin.Context) {
				adminOrderHandler.Ship(c)
			})
			admin.POST("/orders/:id/receive", func(c *gin.Context) {
				adminOrderHandler.Receive(c)
			})
			admin.POST("/orders/:id/cancel", func(c *gin.Context) {
				adminOrderHandler.Cancel(c)
			})
			admin.POST("/orders/:id/refund", func(c *gin.Context) {
				adminOrderHandler.Refund(c)
			})
			admin.POST("/webhooks/dlq/:id/replay", func(c *gin.Context) {
				adminWebhookDLQHandler.Replay(c)
			})
			admin.GET("/webhooks/dlq", func(c *gin.Context) {
				adminWebhookDLQHandler.List(c)
			})
			admin.POST("/oauth/clients", func(c *gin.Context) {
				adminOAuthClientHandler.Create(c)
			})
			admin.GET("/oauth/clients", func(c *gin.Context) {
				adminOAuthClientHandler.List(c)
			})

			adminPaymentHandler := api.NewAdminPaymentHandler(services.Payment)
			admin.GET("/payments", func(c *gin.Context) {
				adminPaymentHandler.List(c)
			})
		}
	}

	r.GET("/.well-known/ucp", func(c *gin.Context) {
		ucpProfileHandler.GetProfile(c)
	})
	r.GET("/.well-known/oauth-authorization-server", func(c *gin.Context) {
		oauthMetadataHandler.WellKnown(c)
	})
	r.POST("/oauth2/token", func(c *gin.Context) {
		oauthTokenHandler.Token(c)
	})
	r.GET("/oauth2/authorize", func(c *gin.Context) {
		oauthAuthorizeHandler.Authorize(c)
	})
	r.GET("/.well-known/oauth-authorization-server", func(c *gin.Context) {
		oauthMetadataHandler.WellKnown(c)
	})
	r.POST("/oauth2/token", func(c *gin.Context) {
		oauthTokenHandler.Token(c)
	})
	r.GET("/oauth2/authorize", func(c *gin.Context) {
		oauthAuthorizeHandler.Authorize(c)
	})

	ucpGroup := r.Group("/ucp/v1")
	{
		ucpGroup.POST("/checkout-sessions", func(c *gin.Context) {
			ucpCheckoutHandler.Create(c)
		})
		ucpGroup.GET("/checkout-sessions/:id", func(c *gin.Context) {
			ucpCheckoutHandler.Get(c)
		})
		ucpGroup.PUT("/checkout-sessions/:id", func(c *gin.Context) {
			ucpCheckoutHandler.Update(c)
		})
		ucpGroup.POST("/checkout-sessions/:id/complete", func(c *gin.Context) {
			ucpCheckoutHandler.Complete(c)
		})
		ucpGroup.DELETE("/checkout-sessions/:id", func(c *gin.Context) {
			ucpCheckoutHandler.Cancel(c)
		})
		ucpGroup.POST("/order-webhooks", func(c *gin.Context) {
			ucpOrderWebhookHandler.Receive(c)
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func buildUCPLinks(links []config.UCPLinkConfig) []ucpmodel.Link {
	result := make([]ucpmodel.Link, 0, len(links))
	for _, link := range links {
		result = append(result, ucpmodel.Link{
			Type:  link.Type,
			URL:   link.URL,
			Title: link.Title,
		})
	}
	return result
}
