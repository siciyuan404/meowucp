package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type auditLogRepository struct {
	db *database.DB
}

func NewAuditLogRepository(db *database.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *domain.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) List(offset, limit int) ([]*domain.AuditLog, error) {
	items := []*domain.AuditLog{}
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *auditLogRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.AuditLog{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
