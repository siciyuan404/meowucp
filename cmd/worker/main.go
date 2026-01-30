package main

import (
	"log"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/internal/ucp/worker"
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

	queueRepo := repository.NewUCPWebhookQueueRepository(db)
	alertRepo := repository.NewUCPWebhookAlertRepository(db)
	processor := worker.NewProcessor(queueRepo, worker.ProcessorConfig{
		BatchSize:   10,
		MaxAttempts: 5,
		BaseDelay:   time.Minute,
	})
	processor.SetAlertSink(worker.NewAlertPolicySink(alertRepo, worker.AlertPolicy{
		MinAttempts:  cfg.UCP.Webhook.AlertMinAttempts,
		DedupeWindow: time.Duration(cfg.UCP.Webhook.AlertDedupeSeconds) * time.Second,
	}))

	sender := worker.NewDeliverySender(cfg.UCP.Webhook.DeliveryURL, time.Duration(cfg.UCP.Webhook.DeliveryTimeoutSec)*time.Second)

	log.Println("Webhook worker started")
	for {
		processed, err := processor.ProcessOnce(func(job *domain.UCPWebhookJob) error {
			return sender.Send(job)
		})
		if err != nil {
			log.Printf("Worker error: %v", err)
		}
		if processed == 0 {
			time.Sleep(2 * time.Second)
		}
	}
}
