package worker

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type AlertPolicy struct {
	MinAttempts  int
	DedupeWindow time.Duration
}

type AlertPolicySink struct {
	repo   repository.UCPWebhookAlertRepository
	policy AlertPolicy
}

func NewAlertPolicySink(repo repository.UCPWebhookAlertRepository, policy AlertPolicy) *AlertPolicySink {
	return &AlertPolicySink{repo: repo, policy: policy}
}

func (s *AlertPolicySink) Notify(alert *domain.UCPWebhookAlert) error {
	if s == nil || s.repo == nil {
		return nil
	}
	minAttempts := s.policy.MinAttempts
	if minAttempts <= 0 {
		minAttempts = 1
	}
	if alert.Attempts < minAttempts {
		return nil
	}

	if s.policy.DedupeWindow > 0 {
		if exists, _ := s.repo.ExistsRecent(alert.EventID, alert.Reason, s.policy.DedupeWindow); exists {
			return nil
		}
	}

	return s.repo.Create(alert)
}
