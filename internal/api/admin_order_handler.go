package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminOrderService interface {
	ListOrders(offset, limit int, filters map[string]interface{}) ([]*domain.Order, int64, error)
	GetOrder(id int64) (*domain.Order, error)
	UpdateOrderStatus(id int64, status string) error
	CancelOrder(id int64, reason string) error
	ShipOrder(id int64, carrier, tracking string) (*domain.Shipment, error)
	ReceiveOrder(id int64) error
}

type AdminOrderHandler struct {
	service AdminOrderService
}

func NewAdminOrderHandler(service AdminOrderService) *AdminOrderHandler {
	return &AdminOrderHandler{service: service}
}

func (h *AdminOrderHandler) List(c *gin.Context) {
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
	if orderNo := c.Query("order_no"); orderNo != "" {
		filters["order_no = ?"] = orderNo
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id = ?"] = userID
	}
	if minAmount := parseFloat(c.Query("amount_min")); minAmount != nil {
		filters["total >= ?"] = *minAmount
	}
	if maxAmount := parseFloat(c.Query("amount_max")); maxAmount != nil {
		filters["total <= ?"] = *maxAmount
	}
	if sku := c.Query("sku"); sku != "" {
		filters["item_sku"] = sku
	}
	if from := parseOrderTime(c.Query("from")); from != nil {
		filters["created_at >= ?"] = *from
	}
	if to := parseOrderTime(c.Query("to")); to != nil {
		filters["created_at <= ?"] = *to
	}

	orders, total, err := h.service.ListOrders(offset, limitInt, filters)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list orders")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func parseOrderTime(value string) *time.Time {
	if value == "" {
		return nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return &parsed
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return &parsed
	}
	return nil
}

func parseFloat(value string) *float64 {
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func (h *AdminOrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid order id")
		return
	}
	order, err := h.service.GetOrder(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "order_not_found", "Order not found")
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *AdminOrderHandler) Ship(c *gin.Context) {
	id, ok := h.parseOrderID(c)
	if !ok {
		return
	}
	order, ok := h.ensureTransitionAllowed(c, id, "shipped")
	if !ok {
		return
	}
	carrier := c.DefaultQuery("carrier", "")
	tracking := c.DefaultQuery("tracking_no", "")
	if _, err := h.service.ShipOrder(order.ID, carrier, tracking); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to ship order")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "shipped"})
}

func (h *AdminOrderHandler) Cancel(c *gin.Context) {
	id, ok := h.parseOrderID(c)
	if !ok {
		return
	}
	_, ok = h.ensureTransitionAllowed(c, id, "cancelled")
	if !ok {
		return
	}
	if err := h.service.CancelOrder(id, "admin_cancel"); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to cancel order")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *AdminOrderHandler) Receive(c *gin.Context) {
	id, ok := h.parseOrderID(c)
	if !ok {
		return
	}
	_, ok = h.ensureTransitionAllowed(c, id, "delivered")
	if !ok {
		return
	}
	if err := h.service.ReceiveOrder(id); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to receive order")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "delivered"})
}

func (h *AdminOrderHandler) Refund(c *gin.Context) {
	h.updateStatus(c, "refunded")
}

func (h *AdminOrderHandler) updateStatus(c *gin.Context, status string) {
	id, ok := h.parseOrderID(c)
	if !ok {
		return
	}
	_, ok = h.ensureTransitionAllowed(c, id, status)
	if !ok {
		return
	}
	if err := h.service.UpdateOrderStatus(id, status); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to update order status")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *AdminOrderHandler) parseOrderID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid order id")
		return 0, false
	}
	return id, true
}

func (h *AdminOrderHandler) ensureTransitionAllowed(c *gin.Context, id int64, status string) (*domain.Order, bool) {
	order, err := h.service.GetOrder(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "order_not_found", "Order not found")
		return nil, false
	}
	if !canTransition(order.Status, status) {
		respondError(c, http.StatusBadRequest, "invalid_status_transition", "Invalid order status transition")
		return nil, false
	}
	return order, true
}

func canTransition(current, next string) bool {
	switch next {
	case "shipped":
		return current == "paid"
	case "cancelled":
		return current == "pending" || current == "paid"
	case "refunded":
		return current == "paid"
	case "delivered":
		return current == "shipped"
	default:
		return false
	}
}
