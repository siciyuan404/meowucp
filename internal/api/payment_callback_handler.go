package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PaymentCallbackPaymentService interface {
	MarkPaymentPaid(orderID int64, transactionID string) error
}

type PaymentCallbackOrderService interface {
	UpdateOrderStatus(id int64, status string) error
}

type PaymentCallbackHandler struct {
	payment PaymentCallbackPaymentService
	orders  PaymentCallbackOrderService
}

func NewPaymentCallbackHandler(payment PaymentCallbackPaymentService, orders PaymentCallbackOrderService) *PaymentCallbackHandler {
	return &PaymentCallbackHandler{payment: payment, orders: orders}
}

type PaymentCallbackRequest struct {
	OrderID       int64  `json:"order_id"`
	TransactionID string `json:"transaction_id"`
}

func (h *PaymentCallbackHandler) Handle(c *gin.Context) {
	var req PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.OrderID <= 0 || strings.TrimSpace(req.TransactionID) == "" {
		respondError(c, http.StatusBadRequest, "missing_required_fields", "Order id and transaction id are required")
		return
	}
	if h.payment == nil || h.orders == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "Payment callback unavailable")
		return
	}
	if err := h.payment.MarkPaymentPaid(req.OrderID, req.TransactionID); err != nil {
		respondError(c, http.StatusInternalServerError, "payment_update_failed", "Failed to update payment")
		return
	}
	if err := h.orders.UpdateOrderStatus(req.OrderID, "paid"); err != nil {
		respondError(c, http.StatusInternalServerError, "order_update_failed", "Failed to update order")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
