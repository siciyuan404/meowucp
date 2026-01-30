//go:build integration
// +build integration

package repository

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/config"
	"github.com/meowucp/pkg/database"
)

func TestProductRepositoryUpdateStockWithDeltaIsAtomic(t *testing.T) {
	db := loadTestDB(t)
	defer db.Close()

	if err := db.AutoMigrate(&domain.Product{}).Error; err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repo := NewProductRepository(db)
	unique := fmt.Sprintf("atomic-%d", time.Now().UnixNano())
	product := &domain.Product{
		Name:          "Atomic Stock",
		Slug:          unique,
		SKU:           unique,
		Price:         10,
		StockQuantity: 5,
		Status:        1,
	}
	if err := repo.Create(product); err != nil {
		t.Fatalf("create product: %v", err)
	}
	defer repo.Delete(product.ID)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	failures := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := repo.UpdateStockWithDelta(product.ID, -1); err != nil {
				mu.Lock()
				failures++
				mu.Unlock()
				return
			}
			mu.Lock()
			successes++
			mu.Unlock()
		}()
	}

	wg.Wait()

	updated, err := repo.FindByID(product.ID)
	if err != nil {
		t.Fatalf("find product: %v", err)
	}
	if successes != 5 {
		t.Fatalf("expected 5 successful decrements, got %d", successes)
	}
	if failures != 5 {
		t.Fatalf("expected 5 failed decrements, got %d", failures)
	}
	if updated.StockQuantity != 0 {
		t.Fatalf("expected stock to be 0, got %d", updated.StockQuantity)
	}
}

func loadTestDB(t *testing.T) *database.DB {
	configPath := os.Getenv("TEST_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	var dbCfg config.DatabaseConfig
	if cfg, err := config.Load(configPath); err == nil {
		dbCfg = cfg.Database
	}

	if host := os.Getenv("TEST_DB_HOST"); host != "" {
		dbCfg.Host = host
	}
	if port := os.Getenv("TEST_DB_PORT"); port != "" {
		if parsed, err := strconv.Atoi(port); err == nil {
			dbCfg.Port = parsed
		}
	}
	if user := os.Getenv("TEST_DB_USER"); user != "" {
		dbCfg.User = user
	}
	if password := os.Getenv("TEST_DB_PASSWORD"); password != "" {
		dbCfg.Password = password
	}
	if name := os.Getenv("TEST_DB_NAME"); name != "" {
		dbCfg.DBName = name
	}
	if sslmode := os.Getenv("TEST_DB_SSLMODE"); sslmode != "" {
		dbCfg.SSLMode = sslmode
	}

	if dbCfg.Host == "" || dbCfg.DBName == "" {
		t.Skip("no database configuration provided for integration test")
	}

	db, err := database.NewDB(
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.DBName,
		dbCfg.Port,
		dbCfg.SSLMode,
		25,
		5,
		300,
	)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	return db
}
