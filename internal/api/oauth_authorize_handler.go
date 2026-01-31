package api

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OAuthAuthorizeHandler struct{}

func NewOAuthAuthorizeHandler() *OAuthAuthorizeHandler {
	return &OAuthAuthorizeHandler{}
}

func (h *OAuthAuthorizeHandler) Authorize(c *gin.Context) {
	responseType := c.Query("response_type")
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")

	if responseType != "code" {
		respondError(c, http.StatusBadRequest, "unsupported_response_type", "Unsupported response type")
		return
	}
	if clientID != oauthClientID || redirectURI == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid client or redirect URI")
		return
	}

	code := "code_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	redirect := redirectURI + "?code=" + url.QueryEscape(code)
	if state != "" {
		redirect += "&state=" + url.QueryEscape(state)
	}

	c.Redirect(http.StatusFound, redirect)
}
