package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeOAuthClientService struct {
	created *domain.OAuthClient
	items   []*domain.OAuthClient
}

func (f *fakeOAuthClientService) CreateClient(client *domain.OAuthClient) error {
	f.created = client
	f.items = append(f.items, client)
	return nil
}

func (f *fakeOAuthClientService) Create(clientID, secret, scopes string) (*domain.OAuthClient, error) {
	client := &domain.OAuthClient{ClientID: clientID, SecretHash: secret, Scopes: scopes, Status: "active"}
	f.created = client
	f.items = append(f.items, client)
	return client, nil
}

func (f *fakeOAuthClientService) ListClients(offset, limit int) ([]*domain.OAuthClient, int64, error) {
	return f.items, int64(len(f.items)), nil
}

func (f *fakeOAuthClientService) List(offset, limit int) ([]*domain.OAuthClient, int64, error) {
	return f.items, int64(len(f.items)), nil
}

func TestAdminCreatesOAuthClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeOAuthClientService{}
	handler := NewAdminOAuthClientHandler(service)

	r := gin.New()
	r.POST("/api/v1/admin/oauth/clients", handler.Create)

	body := `{"client_id":"client_1","secret":"secret","scopes":"checkout"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/oauth/clients", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if service.created == nil || service.created.ClientID != "client_1" {
		t.Fatalf("expected client to be created")
	}
}

func TestAdminListsOAuthClients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeOAuthClientService{items: []*domain.OAuthClient{{ClientID: "client_1"}}}
	handler := NewAdminOAuthClientHandler(service)

	r := gin.New()
	r.GET("/api/v1/admin/oauth/clients", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/oauth/clients", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
