package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type ShippingRateService interface {
	Quote(region string, items []domain.OrderItem) (float64, float64, error)
}

type ShippingRateHandler struct {
	service ShippingRateService
}

func NewShippingRateHandler(service ShippingRateService) *ShippingRateHandler {
	return &ShippingRateHandler{service: service}
}

func (h *ShippingRateHandler) List(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "Shipping service unavailable")
		return
	}
	region := c.Query("region")
	quantity := parseInt(c.DefaultQuery("quantity", "1"))
	if region == "" {
		respondError(c, http.StatusBadRequest, "missing_region", "Region is required")
		return
	}
	if quantity < 1 {
		quantity = 1
	}
	items := []domain.OrderItem{{Quantity: quantity, UnitPrice: 0, TotalPrice: 0}}
	_, shipping, err := h.service.Quote(region, items)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "quote_failed", "Failed to quote shipping")
		return
	}
	_, tax, _ := h.service.Quote(region, items)
	c.JSON(http.StatusOK, gin.H{
		"region":   region,
		"quantity": quantity,
		"shipping": shipping,
		"tax":      tax,
	})
}
