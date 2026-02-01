package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type OAuthTokenRevoker interface {
	Revoke(token string) error
}

type OAuthRevokeHandler struct {
	service OAuthTokenRevoker
}

func NewOAuthRevokeHandler(service OAuthTokenRevoker) *OAuthRevokeHandler {
	return &OAuthRevokeHandler{service: service}
}

func (h *OAuthRevokeHandler) Revoke(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service_unavailable", "OAuth token service unavailable")
		return
	}

	token := strings.TrimSpace(c.PostForm("token"))
	if token == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "Missing token")
		return
	}

	if err := h.service.Revoke(token); err != nil {
		respondError(c, http.StatusInternalServerError, "revoke_failed", "Failed to revoke token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
