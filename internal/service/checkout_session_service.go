package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type CheckoutSessionService struct {
	repo repository.CheckoutSessionRepository
}

func NewCheckoutSessionService(repo repository.CheckoutSessionRepository) *CheckoutSessionService {
	return &CheckoutSessionService{repo: repo}
}

func (s *CheckoutSessionService) Create(session *domain.CheckoutSession) error {
	return s.repo.Create(session)
}

func (s *CheckoutSessionService) Update(session *domain.CheckoutSession) error {
	return s.repo.Update(session)
}

func (s *CheckoutSessionService) GetByID(id string) (*domain.CheckoutSession, error) {
	return s.repo.FindByID(id)
}

func (s *CheckoutSessionService) Delete(id string) error {
	return s.repo.Delete(id)
}
