package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type WebhookAlertLister interface {
	List(offset, limit int) ([]*domain.UCPWebhookAlert, int64, error)
}

type WebhookAlertHandler struct {
	lister WebhookAlertLister
}

func NewWebhookAlertHandler(lister WebhookAlertLister) *WebhookAlertHandler {
	return &WebhookAlertHandler{lister: lister}
}

func (h *WebhookAlertHandler) List(c *gin.Context) {
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
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list webhook alerts")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}
