package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type CouponService interface {
	ValidateCoupon(code string, subtotal float64) (*domain.Coupon, error)
}

type CouponHandler struct {
	service CouponService
}

func NewCouponHandler(service CouponService) *CouponHandler {
	return &CouponHandler{service: service}
}

type couponValidateRequest struct {
	Code     string  `json:"code"`
	Subtotal float64 `json:"subtotal"`
}

func (h *CouponHandler) Validate(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "Coupon service unavailable")
		return
	}
	var req couponValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	code := strings.TrimSpace(req.Code)
	if code == "" {
		respondError(c, http.StatusBadRequest, "invalid_code", "Coupon code required")
		return
	}
	coupon, err := h.service.ValidateCoupon(code, req.Subtotal)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_coupon", "Coupon invalid")
		return
	}
	c.JSON(http.StatusOK, coupon)
}
