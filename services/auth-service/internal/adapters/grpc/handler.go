package grpc

import (
	"context"

	authv1 "github.com/emreEngineering/kervan/gen/go/auth/v1"
	"github.com/emreEngineering/kervan/services/auth-service/internal/application"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	app *application.AuthService
}

func NewAuthHandler(app *application.AuthService) *AuthHandler {
	return &AuthHandler{app: app}
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := h.app.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authv1.RegisterResponse{
		UserId: user.ID,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	token, err := h.app.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, err := h.app.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &authv1.ValidateTokenResponse{
		UserId: userID,
		Valid:  true,
	}, nil
}
