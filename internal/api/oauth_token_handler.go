package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	oauthClientID     = "ucp-client"
	oauthClientSecret = "ucp-secret"
	oauthTokenTTL     = time.Hour
)

type OAuthTokenHandler struct {
	clientRepo repository.OAuthClientRepository
	tokenRepo  repository.OAuthTokenRepository
}

func NewOAuthTokenHandler() *OAuthTokenHandler {
	return &OAuthTokenHandler{}
}

func NewOAuthTokenHandlerWithRepos(clientRepo repository.OAuthClientRepository, tokenRepo repository.OAuthTokenRepository) *OAuthTokenHandler {
	return &OAuthTokenHandler{clientRepo: clientRepo, tokenRepo: tokenRepo}
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
	if h.clientRepo != nil {
		client, err := h.clientRepo.FindByClientID(clientID)
		if err != nil || client == nil {
			respondError(c, http.StatusUnauthorized, "invalid_client", "Invalid client credentials")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
			respondError(c, http.StatusUnauthorized, "invalid_client", "Invalid client credentials")
			return
		}
		if strings.TrimSpace(client.Scopes) == "" {
			respondError(c, http.StatusForbidden, "invalid_scope", "No scopes assigned")
			return
		}
	}

	token, err := generateOAuthToken(clientID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "token_failed", "Failed to issue token")
		return
	}
	if h.tokenRepo != nil {
		_ = h.tokenRepo.Create(&domain.OAuthToken{
			Token:     token,
			ClientID:  clientID,
			Scopes:    oauthCheckoutScope,
			ExpiresAt: time.Now().Add(oauthTokenTTL),
		})
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
