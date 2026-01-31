//go:build tools
// +build tools

package main

import (
	"fmt"
	"log"

	"github.com/meowucp/internal/domain"
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

	var job domain.UCPWebhookJob
	jobErr := db.Order("id desc").First(&job).Error
	if jobErr != nil {
		fmt.Printf("latest_job: error=%v\n", jobErr)
		return
	}

	var alertCount int64
	if err := db.Model(&domain.UCPWebhookAlert{}).Count(&alertCount).Error; err != nil {
		log.Fatalf("Failed to count alerts: %v", err)
	}

	fmt.Printf("latest_job: id=%d event_id=%s status=%s attempts=%d next_retry_at=%s last_error=%s\n",
		job.ID, job.EventID, job.Status, job.Attempts, job.NextRetryAt.Format("2006-01-02 15:04:05"), job.LastError)
	fmt.Printf("alerts: count=%d\n", alertCount)
}
