package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type fakeOAuthTokenService struct{}

func (f *fakeOAuthTokenService) Revoke(token string) error {
	return nil
}

type fakeOAuthClientRepo struct{}

func (f *fakeOAuthClientRepo) Create(client *domain.OAuthClient) error { return nil }
func (f *fakeOAuthClientRepo) FindByClientID(clientID string) (*domain.OAuthClient, error) {
	return &domain.OAuthClient{ClientID: clientID, SecretHash: string(mustHash("secret")), Scopes: "checkout"}, nil
}
func (f *fakeOAuthClientRepo) List(offset, limit int) ([]*domain.OAuthClient, error) { return nil, nil }
func (f *fakeOAuthClientRepo) Count() (int64, error)                                 { return 0, nil }

type fakeOAuthTokenRepo struct{}

func (f *fakeOAuthTokenRepo) Create(token *domain.OAuthToken) error { return nil }
func (f *fakeOAuthTokenRepo) FindByToken(token string) (*domain.OAuthToken, error) {
	return nil, errors.New("not found")
}
func (f *fakeOAuthTokenRepo) Revoke(token string, revokedAt time.Time) error { return nil }

func mustHash(value string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	return hash
}

func TestOAuthTokenExchange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthTokenHandlerWithRepos(nil, nil)
	r := gin.New()
	r.POST("/oauth2/token", handler.Token)

	body := "grant_type=authorization_code&client_id=ucp-client&client_secret=ucp-secret&code=authcode"
	req := httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	accessToken, _ := payload["access_token"].(string)
	if accessToken == "" {
		t.Fatalf("expected access_token to be set")
	}
	if payload["token_type"] != "bearer" {
		t.Fatalf("expected token_type bearer")
	}
	expiresIn, ok := payload["expires_in"].(float64)
	if !ok || expiresIn <= 0 {
		t.Fatalf("expected expires_in to be positive")
	}
}

func TestOAuthTokenScopesRestricted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthTokenHandlerWithRepos(nil, nil)
	r := gin.New()
	r.POST("/oauth2/token", handler.Token)

	clientRepo := &fakeOAuthClientRepo{}
	clientHandler := NewOAuthTokenHandlerWithRepos(clientRepo, &fakeOAuthTokenRepo{})
	r2 := gin.New()
	r2.POST("/oauth2/token", clientHandler.Token)

	body := "grant_type=authorization_code&client_id=client_1&client_secret=wrong-secret&code=authcode"
	req := httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	r2.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestOAuthRevokeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewOAuthRevokeHandler(&fakeOAuthTokenService{})
	r := gin.New()
	r.POST("/oauth2/revoke", handler.Revoke)

	body := "token=token-1"
	req := httptest.NewRequest(http.MethodPost, "/oauth2/revoke", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
