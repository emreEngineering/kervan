package grpc

import (
	"context"
	"time"

	orderv1 "github.com/emreEngineering/kervan/gen/go/order/v1"
	"github.com/emreEngineering/kervan/services/order-service/internal/application"
	"github.com/emreEngineering/kervan/services/order-service/internal/domain"
)

type OrderHandler struct {
	orderv1.UnimplementedOrderServiceServer
	app *application.OrderService
}

func NewOrderHandler(app *application.OrderService) *OrderHandler {
	return &OrderHandler{app: app}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	order, err := h.app.CreateOrder(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &orderv1.CreateOrderResponse{Order: domainOrderToProto(order)}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	order, err := h.app.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &orderv1.GetOrderResponse{Order: domainOrderToProto(order)}, nil
}

func (h *OrderHandler) ReturnOrder(ctx context.Context, req *orderv1.ReturnOrderRequest) (*orderv1.ReturnOrderResponse, error) {
	order, err := h.app.ReturnOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &orderv1.ReturnOrderResponse{Order: domainOrderToProto(order)}, nil
}

func domainOrderToProto(order *domain.Order) *orderv1.Order {
	protoOrder := &orderv1.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
	}
	for _, item := range order.Items {
		protoOrder.Items = append(protoOrder.Items, &orderv1.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	return protoOrder
}
