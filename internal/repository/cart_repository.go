package repository

import (
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/pkg/database"
)

type cartRepository struct {
	db *database.DB
}

func NewCartRepository(db *database.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID int64) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Create(cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) Update(cart *domain.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Cart{}, id).Error
}

func (r *cartRepository) AddItem(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateItem(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) RemoveItem(cartID, productID int64) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&domain.CartItem{}).Error
}

func (r *cartRepository) ClearCart(cartID int64) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}
