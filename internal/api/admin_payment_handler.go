package api

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminPaymentService interface {
	ListPayments(offset, limit int, filters map[string]interface{}) ([]*domain.Payment, int64, error)
}

type AdminPaymentHandler struct {
	service AdminPaymentService
}

func NewAdminPaymentHandler(service AdminPaymentService) *AdminPaymentHandler {
	return &AdminPaymentHandler{service: service}
}

func (h *AdminPaymentHandler) List(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")

	limitInt := parseInt(limit)
	pageInt := parseInt(page)
	if pageInt < 1 {
		pageInt = 1
	}
	if limitInt < 1 {
		limitInt = 20
	}

	offset := (pageInt - 1) * limitInt
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status = ?"] = status
	}
	if method := c.Query("method"); method != "" {
		filters["payment_method = ?"] = method
	}
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id = ?"] = orderID
	}
	if txID := c.Query("transaction_id"); txID != "" {
		filters["transaction_id = ?"] = txID
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id = ?"] = userID
	}
	if currency := c.Query("currency"); currency != "" {
		filters["currency = ?"] = currency
	}
	if minAmount := parseAmountFloat(c.Query("amount_min")); minAmount != nil {
		filters["amount >= ?"] = *minAmount
	}
	if maxAmount := parseAmountFloat(c.Query("amount_max")); maxAmount != nil {
		filters["amount <= ?"] = *maxAmount
	}
	if from := parseOrderTime(c.Query("from")); from != nil {
		filters["created_at >= ?"] = *from
	}
	if to := parseOrderTime(c.Query("to")); to != nil {
		filters["created_at <= ?"] = *to
	}
	items, total, err := h.service.ListPayments(offset, limitInt, filters)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list payments")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func parseAmountFloat(value string) *float64 {
	parsed := parseFloat(value)
	if parsed == nil {
		return nil
	}
	rounded := roundToTwoDecimals(*parsed)
	return &rounded
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100+1e-9) / 100
}
