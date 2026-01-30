package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type fakeInventoryService struct {
	logs []*domain.InventoryLog
}

func (f *fakeInventoryService) AdjustStock(productID int64, change int, notes string) error {
	f.logs = append(f.logs, &domain.InventoryLog{
		ProductID:      productID,
		QuantityChange: change,
		Type:           "adjust",
		Notes:          notes,
		CreatedAt:      time.Now(),
	})
	return nil
}

func (f *fakeInventoryService) ListLogs(productID int64, offset, limit int) ([]*domain.InventoryLog, int64, error) {
	items := make([]*domain.InventoryLog, 0)
	for _, log := range f.logs {
		if log.ProductID == productID {
			items = append(items, log)
		}
	}
	return items, int64(len(items)), nil
}

func TestAdminInventoryAdjust(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeInventoryService{}
	handler := NewAdminInventoryHandler(svc)

	r := gin.New()
	r.POST("/api/v1/admin/inventory/adjust", handler.Adjust)

	payload := map[string]interface{}{"product_id": 1, "quantity_change": 5, "notes": "manual"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/inventory/adjust", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if len(svc.logs) != 1 {
		t.Fatalf("expected one inventory log")
	}
}

func TestAdminInventoryLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &fakeInventoryService{}
	_ = svc.AdjustStock(1, 5, "manual")
	_ = svc.AdjustStock(1, -2, "manual")

	handler := NewAdminInventoryHandler(svc)

	r := gin.New()
	r.GET("/api/v1/admin/inventory/logs", handler.Logs)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/inventory/logs?product_id=1&page=1&limit=10", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "pagination") {
		t.Fatalf("expected pagination in response")
	}
}
