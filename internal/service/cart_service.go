package service

import (
	"errors"

	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/repository"
)

type CartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) GetCart(userID int64) (*domain.Cart, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *CartService) AddToCart(userID, productID int64, quantity int) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		cart = &domain.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return err
		}
	}

	var existingItem *domain.CartItem
	for _, item := range cart.Items {
		if item.ProductID == productID {
			existingItem = &item
			break
		}
	}

	if existingItem != nil {
		existingItem.Quantity += quantity
		return s.cartRepo.UpdateItem(existingItem)
	}

	item := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
	}
	return s.cartRepo.AddItem(item)
}

func (s *CartService) UpdateCartItem(userID, productID int64, quantity int) error {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	var existingItem *domain.CartItem
	for _, item := range cart.Items {
		if item.ProductID == productID {
			existingItem = &item
			break
		}
	}

	if existingItem == nil {
		return errors.New("item not found in cart")
	}

	existingItem.Quantity = quantity
	return s.cartRepo.UpdateItem(existingItem)
}

func (s *CartService) RemoveFromCart(userID, productID int64) error {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	return s.cartRepo.RemoveItem(cart.ID, productID)
}

func (s *CartService) ClearCart(userID int64) error {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	return s.cartRepo.ClearCart(cart.ID)
}
