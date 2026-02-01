package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type idempotencyKeyRepository struct {
	db *database.DB
}

func NewIdempotencyKeyRepository(db *database.DB) IdempotencyKeyRepository {
	return &idempotencyKeyRepository{db: db}
}

func (r *idempotencyKeyRepository) Create(record *domain.IdempotencyKey) error {
	return r.db.Create(record).Error
}

func (r *idempotencyKeyRepository) FindByUserIDAndKey(userID int64, key string) (*domain.IdempotencyKey, error) {
	var record domain.IdempotencyKey
	if err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *idempotencyKeyRepository) Update(record *domain.IdempotencyKey) error {
	return r.db.Save(record).Error
}
