package service

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(category *domain.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CategoryService) GetCategory(id int64) (*domain.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *CategoryService) UpdateCategory(category *domain.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) DeleteCategory(id int64) error {
	return s.categoryRepo.Delete(id)
}

func (s *CategoryService) ListCategories(offset, limit int) ([]*domain.Category, int64, error) {
	categories, err := s.categoryRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.categoryRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (s *CategoryService) GetCategoryTree() ([]*domain.Category, error) {
	return s.categoryRepo.Tree()
}
