package grpc

import (
	"context"

	cartv1 "github.com/emreEngineering/kervan/gen/go/cart/v1"
	"github.com/emreEngineering/kervan/services/cart-service/internal/application"
)

type CartHandler struct {
	cartv1.UnimplementedCartServiceServer
	app *application.CartService
}

func NewCartHandler(app *application.CartService) *CartHandler {
	return &CartHandler{app: app}
}

func (h *CartHandler) GetCart(ctx context.Context, req *cartv1.GetCartRequest) (*cartv1.GetCartResponse, error) {
	cart, err := h.app.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var items []*cartv1.CartItem
	for _, item := range cart.Items {
		items = append(items, &cartv1.CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	return &cartv1.GetCartResponse{Items: items}, nil
}

func (h *CartHandler) AddItem(ctx context.Context, req *cartv1.AddItemRequest) (*cartv1.AddItemResponse, error) {
	err := h.app.AddItem(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &cartv1.AddItemResponse{Success: true}, nil
}

func (h *CartHandler) RemoveItem(ctx context.Context, req *cartv1.RemoveItemRequest) (*cartv1.RemoveItemResponse, error) {
	err := h.app.RemoveItem(ctx, req.UserId, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &cartv1.RemoveItemResponse{Success: true}, nil
}

func (h *CartHandler) ClearCart(ctx context.Context, req *cartv1.ClearCartRequest) (*cartv1.ClearCartResponse, error) {
	err := h.app.ClearCart(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &cartv1.ClearCartResponse{Success: true}, nil
}
