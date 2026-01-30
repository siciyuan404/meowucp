package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type ucpWebhookAuditRepository struct {
	db *database.DB
}

func NewUCPWebhookAuditRepository(db *database.DB) UCPWebhookAuditRepository {
	return &ucpWebhookAuditRepository{db: db}
}

func (r *ucpWebhookAuditRepository) Create(audit *domain.UCPWebhookAudit) error {
	return r.db.Create(audit).Error
}

func (r *ucpWebhookAuditRepository) List(offset, limit int) ([]*domain.UCPWebhookAudit, error) {
	var audits []*domain.UCPWebhookAudit
	err := r.db.Offset(offset).Limit(limit).Order("id desc").Find(&audits).Error
	return audits, err
}

func (r *ucpWebhookAuditRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.UCPWebhookAudit{}).Count(&count).Error
	return count, err
}
