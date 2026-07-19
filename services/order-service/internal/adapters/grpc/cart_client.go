package grpc

import (
	"context"

	cartv1 "github.com/emreEngineering/kervan/gen/go/cart/v1"
	"github.com/emreEngineering/kervan/services/order-service/internal/ports"
	"google.golang.org/grpc"
)

type CartClient struct {
	client cartv1.CartServiceClient
}

func NewCartClient(conn *grpc.ClientConn) *CartClient {
	return &CartClient{client: cartv1.NewCartServiceClient(conn)}
}

func (c *CartClient) GetCart(ctx context.Context, userID int64) ([]ports.CartItemData, error) {
	resp, err := c.client.GetCart(ctx, &cartv1.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	var items []ports.CartItemData
	for _, item := range resp.Items {
		items = append(items, ports.CartItemData{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	return items, nil
}

func (c *CartClient) ClearCart(ctx context.Context, userID int64) error {
	_, err := c.client.ClearCart(ctx, &cartv1.ClearCartRequest{UserId: userID})
	return err
}
