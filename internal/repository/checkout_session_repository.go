package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type checkoutSessionRepository struct {
	db *database.DB
}

func NewCheckoutSessionRepository(db *database.DB) CheckoutSessionRepository {
	return &checkoutSessionRepository{db: db}
}

func (r *checkoutSessionRepository) Create(session *domain.CheckoutSession) error {
	return r.db.Create(session).Error
}

func (r *checkoutSessionRepository) Update(session *domain.CheckoutSession) error {
	return r.db.Save(session).Error
}

func (r *checkoutSessionRepository) FindByID(id string) (*domain.CheckoutSession, error) {
	var session domain.CheckoutSession
	err := r.db.First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *checkoutSessionRepository) Delete(id string) error {
	return r.db.Delete(&domain.CheckoutSession{}, "id = ?", id).Error
}
