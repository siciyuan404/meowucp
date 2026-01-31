package main

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
)

type adminInventoryServiceAdapter struct {
	svc *service.InventoryService
}

func (a adminInventoryServiceAdapter) AdjustStock(productID int64, change int, notes string) error {
	return a.svc.AdjustStock(productID, change, "adjust", "admin", "admin", notes)
}

func (a adminInventoryServiceAdapter) ListLogs(productID int64, offset, limit int) ([]*domain.InventoryLog, int64, error) {
	return a.svc.GetInventoryLogs(productID, offset, limit)
}
