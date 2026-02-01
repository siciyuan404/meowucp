package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/meowucp/internal/domain"
)

type fakeIdempotencyStore struct {
	records map[string]*domain.IdempotencyKey
}

func newFakeIdempotencyStore() *fakeIdempotencyStore {
	return &fakeIdempotencyStore{records: map[string]*domain.IdempotencyKey{}}
}

func (f *fakeIdempotencyStore) Create(record *domain.IdempotencyKey) error {
	key := record.Key
	if _, ok := f.records[key]; ok {
		return gorm.ErrRecordNotFound
	}
	f.records[key] = record
	return nil
}

func (f *fakeIdempotencyStore) FindByUserIDAndKey(userID int64, key string) (*domain.IdempotencyKey, error) {
	record, ok := f.records[key]
	if !ok || record.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return record, nil
}

func (f *fakeIdempotencyStore) Update(record *domain.IdempotencyKey) error {
	f.records[record.Key] = record
	return nil
}

func TestIdempotencyMiddlewareReplaysResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := newFakeIdempotencyStore()
	count := 0

	r := gin.New()
	r.Use(Idempotency(store))
	r.POST("/test", func(c *gin.Context) {
		count++
		c.JSON(http.StatusOK, gin.H{"ok": true, "count": count})
	})

	body := `{"name":"cat"}`
	request := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "key-1")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if count != 1 {
		t.Fatalf("expected handler to run once")
	}

	request2 := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	request2.Header.Set("Content-Type", "application/json")
	request2.Header.Set("Idempotency-Key", "key-1")
	response2 := httptest.NewRecorder()
	r.ServeHTTP(response2, request2)

	if response2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response2.Code)
	}
	if count != 1 {
		t.Fatalf("expected handler to still be 1, got %d", count)
	}
	if response.Body.String() != response2.Body.String() {
		t.Fatalf("expected replayed response to match")
	}
}
