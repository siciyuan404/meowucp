package repository

import (
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type ucpWebhookQueueRepository struct {
	db *database.DB
}

func NewUCPWebhookQueueRepository(db *database.DB) UCPWebhookQueueRepository {
	return &ucpWebhookQueueRepository{db: db}
}

func (r *ucpWebhookQueueRepository) Create(job *domain.UCPWebhookJob) error {
	return r.db.Create(job).Error
}

func (r *ucpWebhookQueueRepository) ListDue(limit int) ([]*domain.UCPWebhookJob, error) {
	var jobs []*domain.UCPWebhookJob
	now := time.Now()
	err := r.db.Where("status IN (?, ?)", "pending", "retrying").
		Where("next_retry_at <= ?", now).
		Order("next_retry_at asc").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}

func (r *ucpWebhookQueueRepository) Update(job *domain.UCPWebhookJob) error {
	return r.db.Save(job).Error
}

func (r *ucpWebhookQueueRepository) List(offset, limit int) ([]*domain.UCPWebhookJob, error) {
	var jobs []*domain.UCPWebhookJob
	err := r.db.Offset(offset).Limit(limit).Order("id desc").Find(&jobs).Error
	return jobs, err
}

func (r *ucpWebhookQueueRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.UCPWebhookJob{}).Count(&count).Error
	return count, err
}

func (r *ucpWebhookQueueRepository) FindByID(id int64) (*domain.UCPWebhookJob, error) {
	var job domain.UCPWebhookJob
	err := r.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}
