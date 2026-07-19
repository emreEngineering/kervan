package ports

import "context"

type ProductClient interface {
	GetProductPrice(ctx context.Context, productID int64) (float64, error)
}
