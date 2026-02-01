package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type PaymentRefundService interface {
	CreateRefund(paymentID int64, amount float64, reason string) (*domain.PaymentRefund, error)
}

type PaymentRefundHandler struct {
	service PaymentRefundService
}

func NewPaymentRefundHandler(service PaymentRefundService) *PaymentRefundHandler {
	return &PaymentRefundHandler{service: service}
}

type PaymentRefundRequest struct {
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}

func (h *PaymentRefundHandler) Create(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "Refund service unavailable")
		return
	}

	paymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || paymentID <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_payment_id", "Invalid payment id")
		return
	}

	var req PaymentRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.Amount <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_amount", "Refund amount required")
		return
	}
	reason := strings.TrimSpace(req.Reason)

	refund, err := h.service.CreateRefund(paymentID, req.Amount, reason)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "refund_failed", "Refund failed")
		return
	}

	c.JSON(http.StatusOK, refund)
}
