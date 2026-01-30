package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type PaymentHandlerService struct {
	repo repository.PaymentHandlerRepository
}

func NewPaymentHandlerService(repo repository.PaymentHandlerRepository) *PaymentHandlerService {
	return &PaymentHandlerService{repo: repo}
}

func (s *PaymentHandlerService) Create(handler *domain.PaymentHandler) error {
	return s.repo.Create(handler)
}

func (s *PaymentHandlerService) Update(handler *domain.PaymentHandler) error {
	return s.repo.Update(handler)
}

func (s *PaymentHandlerService) List() ([]*domain.PaymentHandler, error) {
	return s.repo.List()
}

func (s *PaymentHandlerService) FindByName(name string) (*domain.PaymentHandler, error) {
	return s.repo.FindByName(name)
}
