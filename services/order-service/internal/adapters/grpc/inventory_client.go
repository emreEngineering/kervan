package grpc

import (
	"context"

	inventoryv1 "github.com/emreEngineering/kervan/gen/go/inventory/v1"
	"google.golang.org/grpc"
)

type InventoryClient struct {
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *InventoryClient {
	return &InventoryClient{client: inventoryv1.NewInventoryServiceClient(conn)}
}
func (c *InventoryClient) ReserveStock(ctx context.Context, productID int64, quantity int32) error {
	_, err := c.client.ReserveStock(ctx, &inventoryv1.ReserveStockRequest{ProductId: productID, Quantity: quantity})
	return err
}

func (c *InventoryClient) ReleaseStock(ctx context.Context, productID int64, quantity int32) error {
	_, err := c.client.ReleaseStock(ctx, &inventoryv1.ReleaseStockRequest{ProductId: productID, Quantity: quantity})
	return err
}
