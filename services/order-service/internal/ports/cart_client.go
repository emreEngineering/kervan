package ports

import "context"

type CartItemData struct {
	ProductID int64
	Quantity  int32
}

type CartClient interface {
	GetCart(ctx context.Context, userID int64) ([]CartItemData, error)
	ClearCart(ctx context.Context, userID int64) error
}
