package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminOrderWebhookOrderService interface {
	GetOrder(id int64) (*domain.Order, error)
}

type AdminOrderWebhookQueueService interface {
	EnqueueOrderEvent(order *domain.Order, eventType string) error
	DeliverOrderEvent(order *domain.Order, eventType string, deliveryURL string, timeout time.Duration) error
}

type AdminOrderWebhookConfig struct {
	DeliveryURL string
	Timeout     time.Duration
}

type AdminOrderWebhookHandler struct {
	orderSvc   AdminOrderWebhookOrderService
	webhookSvc AdminOrderWebhookQueueService
	config     AdminOrderWebhookConfig
}

func NewAdminOrderWebhookHandler(orderSvc AdminOrderWebhookOrderService, webhookSvc AdminOrderWebhookQueueService, config AdminOrderWebhookConfig) *AdminOrderWebhookHandler {
	return &AdminOrderWebhookHandler{
		orderSvc:   orderSvc,
		webhookSvc: webhookSvc,
		config:     config,
	}
}

type adminOrderWebhookRequest struct {
	EventType string `json:"event_type"`
	Mode      string `json:"mode"`
}

func (h *AdminOrderWebhookHandler) Trigger(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid order id")
		return
	}

	var req adminOrderWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.EventType == "" {
		respondError(c, http.StatusBadRequest, "missing_event_type", "Event type is required")
		return
	}
	if h.orderSvc == nil || h.webhookSvc == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "Webhook service unavailable")
		return
	}

	order, err := h.orderSvc.GetOrder(orderID)
	if err != nil || order == nil {
		respondError(c, http.StatusNotFound, "order_not_found", "Order not found")
		return
	}

	mode := req.Mode
	if mode == "" {
		mode = "async"
	}

	if mode == "sync" {
		if err := h.webhookSvc.DeliverOrderEvent(order, req.EventType, h.config.DeliveryURL, h.config.Timeout); err != nil {
			respondError(c, http.StatusInternalServerError, "delivery_failed", "Failed to deliver webhook")
			return
		}
	} else {
		if err := h.webhookSvc.EnqueueOrderEvent(order, req.EventType); err != nil {
			respondError(c, http.StatusInternalServerError, "enqueue_failed", "Failed to enqueue webhook")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
