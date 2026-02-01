package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminOAuthClientService interface {
	Create(clientID, secret, scopes string) (*domain.OAuthClient, error)
	List(offset, limit int) ([]*domain.OAuthClient, int64, error)
}

type AdminOAuthClientHandler struct {
	service AdminOAuthClientService
}

func NewAdminOAuthClientHandler(service AdminOAuthClientService) *AdminOAuthClientHandler {
	return &AdminOAuthClientHandler{service: service}
}

type adminOAuthClientRequest struct {
	ClientID string `json:"client_id"`
	Secret   string `json:"secret"`
	Scopes   string `json:"scopes"`
}

func (h *AdminOAuthClientHandler) Create(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "OAuth client service unavailable")
		return
	}
	var req adminOAuthClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	client, err := h.service.Create(req.ClientID, req.Secret, req.Scopes)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "create_failed", "Failed to create client")
		return
	}
	c.JSON(http.StatusOK, client)
}

func (h *AdminOAuthClientHandler) List(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "OAuth client service unavailable")
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
	items, total, err := h.service.List(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list clients")
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
