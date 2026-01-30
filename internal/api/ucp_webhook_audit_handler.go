package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type WebhookAuditLister interface {
	List(offset, limit int) ([]*domain.UCPWebhookAudit, int64, error)
}

type WebhookAuditHandler struct {
	lister WebhookAuditLister
}

func NewWebhookAuditHandler(lister WebhookAuditLister) *WebhookAuditHandler {
	return &WebhookAuditHandler{lister: lister}
}

func (h *WebhookAuditHandler) List(c *gin.Context) {
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
	items, total, err := h.lister.List(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list webhook audits")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audits": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}
