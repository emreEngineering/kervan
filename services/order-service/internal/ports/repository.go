package ports

import (
	"context"

	"github.com/emreEngineering/kervan/services/order-service/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}
