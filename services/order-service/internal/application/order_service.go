package application

import (
	"context"
	"errors"

	"github.com/emreEngineering/kervan/services/order-service/internal/domain"
	"github.com/emreEngineering/kervan/services/order-service/internal/ports"
)

type OrderService struct {
	authRepo      ports.AuthClient
	cartRepo      ports.CartClient
	inventoryRepo ports.InventoryClient
	productRepo   ports.ProductClient
	orderRepo     ports.OrderRepository
}

func NewOrderService(auth ports.AuthClient, cart ports.CartClient, inventory ports.InventoryClient, product ports.ProductClient, order ports.OrderRepository) *OrderService {
	return &OrderService{
		authRepo:      auth,
		cartRepo:      cart,
		inventoryRepo: inventory,
		productRepo:   product,
		orderRepo:     order,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, token string) (*domain.Order, error) {
	userID, err := s.authRepo.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	cartItems, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("sepet boş")
	}

	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		if err := s.inventoryRepo.ReserveStock(ctx, item.ProductID, item.Quantity); err != nil {
			s.releaseReservedStocks(ctx, orderItems)
			return nil, err
		}
		price, err := s.productRepo.GetProductPrice(ctx, item.ProductID)
		if err != nil {
			s.releaseReservedStocks(ctx, orderItems)
			return nil, err
		}
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}
	order := domain.NewOrder(userID, orderItems)
	if err := s.orderRepo.Save(ctx, order); err != nil {
		s.releaseReservedStocks(ctx, orderItems)
		return nil, err
	}
	s.cartRepo.ClearCart(ctx, userID)
	return order, nil
}

func (s *OrderService) releaseReservedStocks(ctx context.Context, orderItems []domain.OrderItem) {
	for _, item := range orderItems {
		s.inventoryRepo.ReleaseStock(ctx, item.ProductID, item.Quantity)
	}
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("sipariş bulunamadı")
	}
	return order, nil
}

func (s *OrderService) ReturnOrder(ctx context.Context, id int64) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("sipariş bulunamadı")
	}

	for _, item := range order.Items {
		s.inventoryRepo.ReleaseStock(ctx, item.ProductID, item.Quantity)
	}
	order.Status = "returned"
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, "returned"); err != nil {
		return nil, err
	}
	return order, nil
}
