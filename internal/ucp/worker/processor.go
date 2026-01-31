package worker

import (
	"time"

	"github.com/meowucp/internal/domain"
)

type QueueStore interface {
	ListDue(limit int) ([]*domain.UCPWebhookJob, error)
	Update(job *domain.UCPWebhookJob) error
}

type ProcessorConfig struct {
	BatchSize   int
	MaxAttempts int
	BaseDelay   time.Duration
}

type Processor struct {
	store     QueueStore
	config    ProcessorConfig
	now       func() time.Time
	alertSink AlertSink
}

type AlertSink interface {
	Notify(alert *domain.UCPWebhookAlert) error
}

func NewProcessor(store QueueStore, config ProcessorConfig) *Processor {
	if config.BatchSize <= 0 {
		config.BatchSize = 10
	}
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = time.Minute
	}
	return &Processor{
		store:  store,
		config: config,
		now:    time.Now,
	}
}

func (p *Processor) SetAlertSink(sink AlertSink) {
	p.alertSink = sink
}

func (p *Processor) ProcessOnce(handler func(job *domain.UCPWebhookJob) error) (int, error) {
	jobs, err := p.store.ListDue(p.config.BatchSize)
	if err != nil {
		return 0, err
	}
	processed := 0
	for _, job := range jobs {
		job.Attempts++
		err := handler(job)
		if err == nil {
			job.Status = "processed"
			job.NextRetryAt = p.now()
		} else {
			job.LastError = err.Error()
			job.LastAttemptAt = p.now()
			if job.Attempts >= p.config.MaxAttempts {
				job.Status = "failed"
			} else {
				job.Status = "retrying"
				job.NextRetryAt = p.now().Add(p.retryDelay(job.Attempts))
			}
			if p.alertSink != nil {
				_ = p.alertSink.Notify(&domain.UCPWebhookAlert{
					EventID:   job.EventID,
					Reason:    "delivery_failed",
					Details:   err.Error(),
					Attempts:  job.Attempts,
					CreatedAt: p.now(),
				})
			}
		}
		if updateErr := p.store.Update(job); updateErr != nil {
			return processed, updateErr
		}
		processed++
	}
	return processed, nil
}

func (p *Processor) retryDelay(attempt int) time.Duration {
	if attempt <= 1 {
		return p.config.BaseDelay
	}
	backoff := time.Duration(1<<uint(attempt-1)) * p.config.BaseDelay
	maxBackoff := time.Duration(1<<uint(p.config.MaxAttempts-1)) * p.config.BaseDelay
	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}
