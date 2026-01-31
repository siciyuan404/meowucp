package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type orderIdempotencyRepository struct {
	db *database.DB
}

func NewOrderIdempotencyRepository(db *database.DB) OrderIdempotencyRepository {
	return &orderIdempotencyRepository{db: db}
}

func (r *orderIdempotencyRepository) Create(record *domain.OrderIdempotency) error {
	return r.db.Create(record).Error
}

func (r *orderIdempotencyRepository) FindByUserIDAndKey(userID int64, key string) (*domain.OrderIdempotency, error) {
	var record domain.OrderIdempotency
	if err := r.db.Where("user_id = ? AND idempotency_key = ?", userID, key).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *orderIdempotencyRepository) Update(record *domain.OrderIdempotency) error {
	return r.db.Save(record).Error
}
