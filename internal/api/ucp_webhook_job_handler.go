package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type WebhookJobAdmin interface {
	List(offset, limit int) ([]*domain.UCPWebhookJob, int64, error)
	RescheduleNow(id int64) error
}

type WebhookJobHandler struct {
	admin WebhookJobAdmin
}

func NewWebhookJobHandler(admin WebhookJobAdmin) *WebhookJobHandler {
	return &WebhookJobHandler{admin: admin}
}

func (h *WebhookJobHandler) List(c *gin.Context) {
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
	items, total, err := h.admin.List(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list webhook jobs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func (h *WebhookJobHandler) Retry(c *gin.Context) {
	idParam := c.Param("id")
	jobID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid job id")
		return
	}

	if err := h.admin.RescheduleNow(jobID); err != nil {
		respondError(c, http.StatusInternalServerError, "reschedule_failed", "Failed to reschedule job")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
