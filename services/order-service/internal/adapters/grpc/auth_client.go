package grpc

import (
	"context"

	authv1 "github.com/emreEngineering/kervan/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client authv1.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{client: authv1.NewAuthServiceClient(conn)}
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (int64, error) {
	resp, err := c.client.ValidateToken(ctx, &authv1.ValidateTokenRequest{Token: token})
	if err != nil {
		return 0, err
	}
	return resp.UserId, nil
}
