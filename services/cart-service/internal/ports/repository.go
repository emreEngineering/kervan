package ports

import (
	"context"

	"github.com/emreEngineering/kervan/services/cart-service/internal/domain"
)

type CartRepository interface {
	GetCart(ctx context.Context, userID int64) (*domain.Cart, error)
	SaveCart(ctx context.Context, cart *domain.Cart) error
	ClearCart(ctx context.Context, userID int64) error
}
