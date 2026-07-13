package ports

import (
	"context"

	"github.com/emreEngineering/kervan/services/inventory-service/internal/domain"
)

type StockRepository interface {
	GetStock(ctx context.Context, productID int64) (*domain.Stock, error)
	Save(ctx context.Context, stock *domain.Stock) error
	Update(ctx context.Context, domain *domain.Stock) error
}
