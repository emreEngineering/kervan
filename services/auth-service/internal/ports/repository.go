package ports

import (
	"context"

	"github.com/emreEngineering/kervan/services/auth-service/internal/domain"
)

// Repository, verilerin saklandığı yere erişen katmandır.
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
