package ports

import (
	"context"

	"github.com/emreEngineering/kervan/service/auth-service/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
