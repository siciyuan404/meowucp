package seed

import (
	"encoding/json"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

const nowPaymentsName = "com.nowpayments"

type NowPaymentsSeedConfig struct {
	Spec         string
	ConfigSchema string
	APIBase      string
	Environment  string
}

func SeedNowPayments(repo repository.PaymentHandlerRepository, cfg NowPaymentsSeedConfig) error {
	if repo == nil {
		return nil
	}

	if existing, err := repo.FindByName(nowPaymentsName); err == nil && existing != nil {
		return nil
	}

	configValue := map[string]string{
		"api_base":    cfg.APIBase,
		"environment": cfg.Environment,
	}
	configJSON, err := json.Marshal(configValue)
	if err != nil {
		return err
	}

	handler := &domain.PaymentHandler{
		Name:         nowPaymentsName,
		Version:      "2026-01-11",
		Spec:         cfg.Spec,
		ConfigSchema: cfg.ConfigSchema,
		Config:       string(configJSON),
	}

	return repo.Create(handler)
}
