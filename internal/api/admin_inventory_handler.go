package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminInventoryService interface {
	AdjustStock(productID int64, change int, notes string) error
	ListLogs(productID int64, offset, limit int) ([]*domain.InventoryLog, int64, error)
}

type AdminInventoryHandler struct {
	service AdminInventoryService
}

func NewAdminInventoryHandler(service AdminInventoryService) *AdminInventoryHandler {
	return &AdminInventoryHandler{service: service}
}

type AdminInventoryAdjustRequest struct {
	ProductID      int64  `json:"product_id"`
	QuantityChange int    `json:"quantity_change"`
	Notes          string `json:"notes"`
}

func (h *AdminInventoryHandler) Adjust(c *gin.Context) {
	var req AdminInventoryAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.ProductID <= 0 || req.QuantityChange == 0 {
		respondError(c, http.StatusBadRequest, "missing_required_fields", "Product id and quantity change are required")
		return
	}
	if err := h.service.AdjustStock(req.ProductID, req.QuantityChange, req.Notes); err != nil {
		respondError(c, http.StatusInternalServerError, "adjust_failed", "Failed to adjust inventory")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AdminInventoryHandler) Logs(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64)
	if err != nil || productID <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_product_id", "Invalid product id")
		return
	}

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
	items, total, err := h.service.ListLogs(productID, offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list inventory logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}
