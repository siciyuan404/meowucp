package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
)

type OrderService interface {
	CreateOrder(userID int64, idempotencyKey string, shippingAddress, billingAddress, notes string, paymentMethod string) (*domain.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type OrderCreateRequest struct {
	UserID          int64  `json:"user_id" binding:"required"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
	BillingAddress  string `json:"billing_address" binding:"required"`
	Notes           string `json:"notes"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.UserID <= 0 {
		respondError(c, http.StatusBadRequest, "missing_required_fields", "User id is required")
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	order, err := h.service.CreateOrder(req.UserID, idempotencyKey, req.ShippingAddress, req.BillingAddress, req.Notes, req.PaymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderIdempotencyConflict):
			respondError(c, http.StatusConflict, "idempotency_conflict", "Order idempotency conflict")
		case err.Error() == "cart not found":
			respondError(c, http.StatusNotFound, "cart_not_found", "Cart not found")
		case err.Error() == "cart is empty":
			respondError(c, http.StatusBadRequest, "cart_empty", "Cart is empty")
		case err.Error() == "insufficient stock":
			respondError(c, http.StatusConflict, "insufficient_stock", "Insufficient stock")
		case err.Error() == "product not found":
			respondError(c, http.StatusNotFound, "product_not_found", "Product not found")
		default:
			respondError(c, http.StatusInternalServerError, "create_failed", "Failed to create order")
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}
