package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type orderStatusLogRepository struct {
	db *database.DB
}

func NewOrderStatusLogRepository(db *database.DB) OrderStatusLogRepository {
	return &orderStatusLogRepository{db: db}
}

func (r *orderStatusLogRepository) Create(log *domain.OrderStatusLog) error {
	return r.db.Create(log).Error
}
