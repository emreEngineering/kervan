package application

import (
	"context"
	"testing"

	"github.com/emreEngineering/kervan/service/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type fakeUserRepo struct {
	users map[string]*domain.User
}

func newFakeUserRepository() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *fakeUserRepo) Save(ctx context.Context, user *domain.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func TestLogin_Success(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	// Önce kullanıcıyı kaydet
	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("kayıt başarısız: %v", err)
	}

	// Doğru şifre ile giriş dene
	token, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("giriş başarısız: %v", err)
	}
	if token == "" {
		t.Fatal("token boş olamaz")
	}
}

func TestRegister_Success(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	user, err := svc.Register(context.Background(), "test@example.com", "password123")

	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if user == nil {
		t.Fatal("user nil olamaz")
	}
	if user.Email != "test@example.com" {
		t.Errorf("email = %s, beklenen test@example.com", user.Email)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("ilk kayıta hata beklenmiyordu: %v", err)
	}

	_, err = svc.Register(context.Background(), "test@example.com", "password456")
	if err == nil {
		t.Fatal("aynı email ile ikinci kayıtta hata bekleniyordu")
	}
}

func TestRegister_EmptyEmail(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Register(context.Background(), "", "password123")

	if err == nil {
		t.Fatal("boş email için hata bekleniyordu")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Login(context.Background(), "missing@example.com", "password123")

	if err == nil {
		t.Fatal("olamayan email için hata bekleniyordu")
	}

}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("kayıt başarısız: %v", err)
	}
	_, err = svc.Login(context.Background(), "test@example.com", "wrong-password")
	if err == nil {
		t.Fatal("yanlış şifre için hata bekleniyordu")
	}

}
func TestValidateToken_Success(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("kayıt başarısız: %v", err)
	}

	token, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("giriş başarısız: %v", err)
	}
	t.Logf("üretilen token: %s", token)

	parsed, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err == nil {
		t.Logf("Token Valid: %v, Claims: %+v", parsed.Valid, parsed.Claims)
	}

	userID, err := svc.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("token doğrulama başarısız: %v", err)
	}
	if err != nil {
		t.Fatalf("token doğrulama başarısız: %v", err)
	}
	// userID 0 olabilir çünkü fake repo ID atamaz
	t.Logf("doğrulanan userID: %d", userID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewAuthService(repo, "test-secret")

	_, err := svc.ValidateToken(context.Background(), "bu-geçersiz-bir-token")

	if err == nil {
		t.Fatal("geçersiz token için hata bekleniyordu")
	}
}
