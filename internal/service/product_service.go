package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
	"github.com/meowucp/pkg/redis"
)

type ProductService struct {
	productRepo   repository.ProductRepository
	inventoryRepo repository.InventoryRepository
	redis         *redis.Client
	cache         ProductListCache
}

func NewProductService(productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, redis *redis.Client) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		redis:         redis,
	}
}

type ProductListCache interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

func NewProductServiceWithCache(productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, cache ProductListCache) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		cache:         cache,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) GetProduct(id int64) (*domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	cached, err := s.redis.Get(context.Background(), cacheKey)
	if err == nil && cached != "" {
		var product domain.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			go s.productRepo.IncrementViews(id)
			return &product, nil
		}
	}

	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	go s.productRepo.IncrementViews(id)

	if data, err := json.Marshal(product); err == nil {
		s.redis.Set(context.Background(), cacheKey, string(data), time.Hour)
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(product *domain.Product) error {
	err := s.productRepo.Update(product)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", product.ID)
	s.redis.Del(context.Background(), cacheKey)

	return nil
}

func (s *ProductService) DeleteProduct(id int64) error {
	cacheKey := fmt.Sprintf("product:%d", id)
	s.redis.Del(context.Background(), cacheKey)

	return s.productRepo.Delete(id)
}

func (s *ProductService) ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error) {
	if s.cache != nil {
		cacheKey := buildProductListCacheKey(offset, limit, filters)
		if cached, ok := s.cache.Get(cacheKey); ok {
			var payload productListCachePayload
			if err := json.Unmarshal([]byte(cached), &payload); err == nil {
				return payload.Items, payload.Total, nil
			}
		}
	}
	products, err := s.productRepo.List(offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.productRepo.Count(filters)
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		cacheKey := buildProductListCacheKey(offset, limit, filters)
		payload := productListCachePayload{Items: products, Total: count}
		if data, err := json.Marshal(payload); err == nil {
			s.cache.Set(cacheKey, string(data))
		}
	}

	return products, count, nil
}

type productListCachePayload struct {
	Items []*domain.Product `json:"items"`
	Total int64             `json:"total"`
}

func buildProductListCacheKey(offset, limit int, filters map[string]interface{}) string {
	filterData, _ := json.Marshal(filters)
	return fmt.Sprintf("products:list:%d:%d:%s", offset, limit, string(filterData))
}

func (s *ProductService) SearchProducts(query string, offset, limit int) ([]*domain.Product, int64, error) {
	products, err := s.productRepo.Search(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.productRepo.SearchCount(query)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}
