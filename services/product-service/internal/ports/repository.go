package ports

import (
	"context"

	"github.com/emreEngineering/kervan/services/product-service/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}
