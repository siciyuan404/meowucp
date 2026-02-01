package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type i18nStringRepository struct {
	db *database.DB
}

func NewI18nStringRepository(db *database.DB) I18nStringRepository {
	return &i18nStringRepository{db: db}
}

func (r *i18nStringRepository) FindByKeyAndLocale(key, locale string) (*domain.I18nString, error) {
	var str domain.I18nString
	if err := r.db.Where("key = ? AND locale = ?", key, locale).First(&str).Error; err != nil {
		return nil, err
	}
	return &str, nil
}

func (r *i18nStringRepository) Create(str *domain.I18nString) error {
	return r.db.Create(str).Error
}
