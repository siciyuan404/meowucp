package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type categoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *categoryRepository) FindByID(id int64) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) List(offset, limit int) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Order("sort_order ASC").Offset(offset).Limit(limit).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Count(&count).Error
	return count, err
}

func (r *categoryRepository) Tree() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}
