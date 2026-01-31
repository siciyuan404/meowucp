//go:build tools
// +build tools

package main

import (
	"log"

	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/seed"
	"github.com/meowucp/pkg/config"
	"github.com/meowucp/pkg/database"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPaymentHandlerRepository(db)
	seedConfig := seed.NowPaymentsSeedConfig{
		Spec:         "https://nowpayments.io",
		ConfigSchema: "https://nowpayments.io",
		APIBase:      "https://api.nowpayments.io",
		Environment:  "test",
	}

	if err := seed.SeedNowPayments(repo, seedConfig); err != nil {
		log.Fatalf("Failed to seed payment handlers: %v", err)
	}

	log.Println("Payment handlers seeded successfully")
}
