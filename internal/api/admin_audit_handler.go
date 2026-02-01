package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AuditLogService interface {
	List(offset, limit int) ([]*domain.AuditLog, int64, error)
}

type AdminAuditHandler struct {
	service AuditLogService
}

func NewAdminAuditHandler(service AuditLogService) *AdminAuditHandler {
	return &AdminAuditHandler{service: service}
}

func (h *AdminAuditHandler) List(c *gin.Context) {
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

	items, total, err := h.service.List(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed_to_list_audit_logs", "failed to list audit logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
		"data": items,
	})
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
