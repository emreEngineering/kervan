package application

import (
	"context"
	"errors"
	"testing"

	"github.com/emreEngineering/kervan/services/order-service/internal/domain"
	"github.com/emreEngineering/kervan/services/order-service/internal/ports"
)

type fakeAuthClient struct {
	userID int64
	err    error
}

func (f *fakeAuthClient) ValidateToken(_ context.Context, _ string) (int64, error) {
	return f.userID, f.err
}

type fakeCartClient struct {
	items   []ports.CartItemData
	err     error
	cleared bool
}

func (f *fakeCartClient) GetCart(_ context.Context, _ int64) ([]ports.CartItemData, error) {
	return f.items, f.err
}

func (f *fakeCartClient) ClearCart(_ context.Context, _ int64) error {
	f.cleared = true
	return nil
}

type fakeInventoryClient struct {
	reserved   []ports.CartItemData
	released   []ports.CartItemData
	reserveErr map[int64]error
}

func (f *fakeInventoryClient) ReserveStock(_ context.Context, productID int64, quantity int32) error {
	if f.reserveErr != nil {
		if err, ok := f.reserveErr[productID]; ok {
			return err
		}
	}
	f.reserved = append(f.reserved, ports.CartItemData{ProductID: productID, Quantity: quantity})
	return nil
}

func (f *fakeInventoryClient) ReleaseStock(_ context.Context, productID int64, quantity int32) error {
	f.released = append(f.released, ports.CartItemData{ProductID: productID, Quantity: quantity})
	return nil
}

type fakeProductClient struct {
	prices map[int64]float64
	err    error
}

func (f *fakeProductClient) GetProductPrice(_ context.Context, productID int64) (float64, error) {
	if f.err != nil {
		return 0, f.err
	}

	price, ok := f.prices[productID]
	if !ok {
		return 0, errors.New("fiyat bulunamadı")
	}
	return price, nil
}

type fakeOrderRepo struct {
	records []domain.Order
}

func (f *fakeOrderRepo) Save(_ context.Context, order *domain.Order) error {
	order.ID = int64(len(f.records) + 1)
	f.records = append(f.records, *order)
	return nil
}

func (f *fakeOrderRepo) FindByID(_ context.Context, id int64) (*domain.Order, error) {
	for _, o := range f.records {
		if o.ID == id {
			return &o, nil
		}
	}
	return nil, nil
}

func (f *fakeOrderRepo) UpdateStatus(_ context.Context, id int64, status domain.OrderStatus) error {
	for i, o := range f.records {
		if o.ID == id {
			f.records[i].Status = status
			return nil
		}
	}
	return errors.New("bulunamadı")
}

// Testler

func TestCreateOrder_Succes(t *testing.T) {
	auth := &fakeAuthClient{userID: 1}
	cart := &fakeCartClient{
		items: []ports.CartItemData{
			{ProductID: 100, Quantity: 2},
			{ProductID: 200, Quantity: 1},
		},
	}

	inventory := &fakeInventoryClient{}
	product := &fakeProductClient{prices: map[int64]float64{100: 50, 200: 100}}
	orderRepo := &fakeOrderRepo{}

	svc := NewOrderService(auth, cart, inventory, product, orderRepo)
	order, err := svc.CreateOrder(context.Background(), "valid-token")

	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if order.TotalAmount != 200 {
		t.Fatalf("beklenen TotalAmount=200, gerçek=%f", order.TotalAmount)
	}
	if order.Status != domain.OrderStatusPending {
		t.Fatalf("beklenen status=pending, gerçek=%s", order.Status)
	}
	if !cart.cleared {
		t.Fatal("sepet temizlenemedi")
	}
	if len(inventory.reserved) != 2 {
		t.Fatalf("beklenen 2 rezervasyon, gerçek %d", len(inventory.reserved))
	}
}

func TestCreateOrder_InvalidToken(t *testing.T) {
	auth := &fakeAuthClient{err: errors.New("geçersiz token")}
	cart := &fakeCartClient{}
	inventory := &fakeInventoryClient{}
	product := &fakeProductClient{}
	orderRepo := &fakeOrderRepo{}

	svc := NewOrderService(auth, cart, inventory, product, orderRepo)
	_, err := svc.CreateOrder(context.Background(), "invalid-token")
	if err == nil {
		t.Fatal("token hatası bekleniyordu")
	}
	if cart.cleared {
		t.Fatal("sepet temizlenmemeliydi")
	}
}

func TestCreateOrder_EmptyCart(t *testing.T) {
	auth := &fakeAuthClient{userID: 1}
	cart := &fakeCartClient{items: []ports.CartItemData{}}
	inventory := &fakeInventoryClient{}
	product := &fakeProductClient{}
	orderRepo := &fakeOrderRepo{}

	svc := NewOrderService(auth, cart, inventory, product, orderRepo)
	_, err := svc.CreateOrder(context.Background(), "valid-token")
	if err == nil {
		t.Fatal("boş sepet hatası bekleniyordu")
	}
}

func TestCreateOrder_ReserveStockFails(t *testing.T) {
	auth := &fakeAuthClient{userID: 1}
	cart := &fakeCartClient{
		items: []ports.CartItemData{
			{ProductID: 100, Quantity: 2},
			{ProductID: 200, Quantity: 1},
		},
	}
	inventory := &fakeInventoryClient{
		reserveErr: map[int64]error{200: errors.New("yetersiz stok")},
	}
	product := &fakeProductClient{prices: map[int64]float64{100: 50, 200: 100}}
	orderRepo := &fakeOrderRepo{}

	svc := NewOrderService(auth, cart, inventory, product, orderRepo)
	_, err := svc.CreateOrder(context.Background(), "valid-token")
	if err == nil {
		t.Fatal("stok hatası bekleniyordu")
	}

	if len(inventory.released) != 1 || inventory.released[0].ProductID != 100 {
		t.Fatalf("compensating action çalımadı, released: %+v", inventory.released)
	}
}
