package ports

import "context"

type InventoryClient interface {
	ReserveStock(ctx context.Context, productID int64, quantity int32) error
	ReleaseStock(ctx context.Context, productID int64, quantity int32) error
}
