package application

import (
	"context"
	"testing"

	"github.com/emreEngineering/kervan/services/cart-service/internal/domain"
)

type fakeCartRepo struct {
	carts map[int64]*domain.Cart
}

func (f *fakeCartRepo) GetCart(_ context.Context, userID int64) (*domain.Cart, error) {
	cart, ok := f.carts[userID]
	if !ok {
		return nil, nil
	}
	return cart, nil
}

func (f *fakeCartRepo) SaveCart(_ context.Context, cart *domain.Cart) error {
	f.carts[cart.UserID] = cart
	return nil
}

func (f *fakeCartRepo) ClearCart(_ context.Context, userID int64) error {
	delete(f.carts, userID)
	return nil
}

func TestGetCart_NewUser(t *testing.T) {
	repo := &fakeCartRepo{carts: make(map[int64]*domain.Cart)}
	svc := NewCartService(repo)
	cart, err := svc.GetCart(context.Background(), 1)

	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	if cart.UserID != 1 || len(cart.Items) != 0 {
		t.Fatal("beklenen userID=1, boş sepet")
	}
}

func TestAddItem_Success(t *testing.T) {
	repo := &fakeCartRepo{carts: make(map[int64]*domain.Cart)}
	svc := NewCartService(repo)

	err := svc.AddItem(context.Background(), 1, 100, 2)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	cart, _ := svc.GetCart(context.Background(), 1)
	if len(cart.Items) != 1 || cart.Items[0].Quantity != 2 {
		t.Fatalf("beklenen 1 ürün 2 adet")
	}

}

func TestAddItem_ExistingProduct(t *testing.T) {
	repo := &fakeCartRepo{carts: make(map[int64]*domain.Cart)}
	svc := NewCartService(repo)

	svc.AddItem(context.Background(), 1, 100, 2)
	svc.AddItem(context.Background(), 1, 100, 3)

	cart, _ := svc.GetCart(context.Background(), 1)

	if len(cart.Items) != 1 || cart.Items[0].Quantity != 5 {
		t.Fatalf("beklenen ürün 1 adet 5")
	}

}

func TestRemoveItem_Success(t *testing.T) {
	repo := &fakeCartRepo{carts: make(map[int64]*domain.Cart)}
	svc := NewCartService(repo)

	svc.AddItem(context.Background(), 1, 100, 2)
	svc.AddItem(context.Background(), 1, 200, 1)
	err := svc.RemoveItem(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	cart, _ := svc.GetCart(context.Background(), 1)
	if len(cart.Items) != 1 || cart.Items[0].ProductID != 200 {
		t.Fatalf("beklenen 1 ürün(200)")
	}
}

func TestClearCart(t *testing.T) {
	repo := &fakeCartRepo{carts: make(map[int64]*domain.Cart)}
	svc := NewCartService(repo)

	svc.AddItem(context.Background(), 1, 100, 2)
	err := svc.ClearCart(context.Background(), 1)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	cart, _ := svc.GetCart(context.Background(), 1)
	if len(cart.Items) != 0 {
		t.Fatalf("sepet boş olmalı")
	}
}
