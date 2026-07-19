package ports

import "context"

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (int64, error)
}
