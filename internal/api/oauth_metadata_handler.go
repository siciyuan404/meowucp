package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const oauthCheckoutScope = "ucp:scopes:checkout_session"

type OAuthMetadataHandler struct{}

func NewOAuthMetadataHandler() *OAuthMetadataHandler {
	return &OAuthMetadataHandler{}
}

func (h *OAuthMetadataHandler) WellKnown(c *gin.Context) {
	baseURL := resolveOAuthBaseURL(c)

	c.JSON(http.StatusOK, gin.H{
		"issuer":                                baseURL,
		"authorization_endpoint":                baseURL + "/oauth2/authorize",
		"token_endpoint":                        baseURL + "/oauth2/token",
		"scopes_supported":                      []string{oauthCheckoutScope},
		"response_types_supported":              []string{"code"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_post"},
	})
}

func resolveOAuthBaseURL(c *gin.Context) string {
	proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	host := c.Request.Host
	return proto + "://" + host
}
