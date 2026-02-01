package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminWebhookDLQService interface {
	ListDLQ(offset, limit int) ([]*domain.WebhookDLQ, int64, error)
	ReplayDLQ(id int64) error
}

type AdminWebhookDLQHandler struct {
	service AdminWebhookDLQService
}

func NewAdminWebhookDLQHandler(service AdminWebhookDLQService) *AdminWebhookDLQHandler {
	return &AdminWebhookDLQHandler{service: service}
}

func (h *AdminWebhookDLQHandler) List(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "DLQ service unavailable")
		return
	}
	limitInt := parseInt(c.DefaultQuery("limit", "20"))
	pageInt := parseInt(c.DefaultQuery("page", "1"))
	if pageInt < 1 {
		pageInt = 1
	}
	if limitInt < 1 {
		limitInt = 20
	}

	offset := (pageInt - 1) * limitInt
	items, total, err := h.service.ListDLQ(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list DLQ")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func (h *AdminWebhookDLQHandler) Replay(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "DLQ service unavailable")
		return
	}
	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || jobID <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid DLQ id")
		return
	}
	if err := h.service.ReplayDLQ(jobID); err != nil {
		respondError(c, http.StatusInternalServerError, "replay_failed", "Failed to replay DLQ")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
