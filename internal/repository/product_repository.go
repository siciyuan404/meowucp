package repository

import (
	"errors"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type productRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) FindByID(id int64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByIDs(ids []int64) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}
	var products []*domain.Product
	err := r.db.Where("id IN ?", ids).Find(&products).Error
	return products, err
}

func (r *productRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.Product, error) {
	var products []*domain.Product
	query := r.db.Offset(offset).Limit(limit)

	for key, value := range filters {
		query = query.Where(key, value)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *productRepository) Count(filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&domain.Product{})

	for key, value := range filters {
		query = query.Where(key, value)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *productRepository) Search(query string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

func (r *productRepository) SearchCount(query string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Product{}).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error
	return count, err
}

func (r *productRepository) UpdateStock(id int64, quantity int) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Update("stock_quantity", quantity).Error
}

func (r *productRepository) UpdateStockWithDelta(id int64, delta int) error {
	result := r.db.Exec(
		"UPDATE products SET stock_quantity = stock_quantity + ? WHERE id = ? AND stock_quantity + ? >= 0",
		delta,
		id,
		delta,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepository) IncrementViews(id int64) error {
	return r.db.Exec("UPDATE products SET views = views + 1 WHERE id = ?", id).Error
}

func (r *productRepository) IncrementSales(id int64, quantity int) error {
	return r.db.Exec("UPDATE products SET sales = sales + ? WHERE id = ?", quantity, id).Error
}
