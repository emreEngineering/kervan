package application

import (
	"context"
	"errors"
	"time"

	"github.com/emreEngineering/kervan/services/auth-service/internal/domain"
	"github.com/emreEngineering/kervan/services/auth-service/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      ports.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo ports.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email ve şifre boş olamaz")
	}

	existing, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("bu email zaten kayıtlı")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Boş kontrol
	if email == "" || password == "" {
		return "", errors.New("email ve şifre zorunludur")
	}

	// 2. Kullanıcıyı email ile bul
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("email veya şifre hatalı")
	}

	// 3. Şifreyi karşılaştır
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("email veya şifre hatalı")
	}

	// 4. JWT token üret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * 365 * 10 * time.Hour).Unix()})
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("beklenmeyen imza yöntemi")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("geçersiz token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("token içinde user_id bulunamadı")
	}

	return int64(userIDFloat), nil
}
