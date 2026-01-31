package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	oauthClientID     = "ucp-client"
	oauthClientSecret = "ucp-secret"
	oauthTokenTTL     = time.Hour
)

type OAuthTokenHandler struct{}

func NewOAuthTokenHandler() *OAuthTokenHandler {
	return &OAuthTokenHandler{}
}

func (h *OAuthTokenHandler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	code := c.PostForm("code")

	if grantType != "authorization_code" {
		respondError(c, http.StatusBadRequest, "unsupported_grant_type", "Unsupported grant type")
		return
	}
	if clientID == "" || clientSecret == "" || code == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "Missing required fields")
		return
	}
	if clientID != oauthClientID || clientSecret != oauthClientSecret {
		respondError(c, http.StatusUnauthorized, "invalid_client", "Invalid client credentials")
		return
	}

	token, err := generateOAuthToken(clientID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "token_failed", "Failed to issue token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   int(oauthTokenTTL.Seconds()),
		"scope":        oauthCheckoutScope,
	})
}

func generateOAuthToken(clientID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   clientID,
		"scope": oauthCheckoutScope,
		"typ":   "oauth",
		"exp":   time.Now().Add(oauthTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
