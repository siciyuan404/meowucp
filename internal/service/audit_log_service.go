package service

import (
	"errors"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type AuditLogService struct {
	repo repository.AuditLogRepository
}

func NewAuditLogService(repo repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) Record(actor, action, target string, payload string) error {
	if s == nil || s.repo == nil {
		return errors.New("audit_log_repo_unavailable")
	}
	var payloadPtr *string
	if payload != "" {
		payloadPtr = &payload
	}
	return s.repo.Create(&domain.AuditLog{
		Actor:     actor,
		Action:    action,
		Target:    target,
		Payload:   payloadPtr,
		CreatedAt: time.Now(),
	})
}

func (s *AuditLogService) List(offset, limit int) ([]*domain.AuditLog, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, errors.New("audit_log_repo_unavailable")
	}
	items, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
