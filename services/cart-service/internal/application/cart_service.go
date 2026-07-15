package application

import (
	"context"

	"github.com/emreEngineering/kervan/services/cart-service/internal/domain"
	"github.com/emreEngineering/kervan/services/cart-service/internal/ports"
)

type CartService struct {
	repo ports.CartRepository
}

func NewCartService(repo ports.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	cart, err := s.repo.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}
	if cart == nil {
		return domain.NewCart(userID), nil
	}
	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, userID, productID int64, quantity int32) error {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = domain.NewCart(userID)
	}

	if err := cart.AddItem(productID, quantity); err != nil {
		return err
	}
	return s.repo.SaveCart(ctx, cart)
}

func (s *CartService) RemoveItem(ctx context.Context, userID, productID int64) error {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = domain.NewCart(userID)
	}

	if err := cart.RemoveItem(productID); err != nil {
		return err
	}
	return s.repo.SaveCart(ctx, cart)
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	return s.repo.ClearCart(ctx, userID)
}
