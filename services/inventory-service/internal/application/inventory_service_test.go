package application

import (
	"context"
	"testing"

	"github.com/emreEngineering/kervan/services/inventory-service/internal/domain"
)

type fakeStockRepo struct {
	stocks map[int64]*domain.Stock
}

func (f *fakeStockRepo) GetStock(_ context.Context, productID int64) (*domain.Stock, error) {
	s, ok := f.stocks[productID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (f *fakeStockRepo) Save(_ context.Context, stock *domain.Stock) error {
	f.stocks[stock.ProductID] = stock
	return nil
}

func (f *fakeStockRepo) Update(_ context.Context, stock *domain.Stock) error {
	f.stocks[stock.ProductID] = stock
	return nil
}

func TestSetStock_Success(t *testing.T) {
	repo := &fakeStockRepo{
		stocks: make(map[int64]*domain.Stock),
	}
	svc := NewInventoryService(repo)

	stock, err := svc.SetStock(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if stock.Available != 10 {
		t.Fatalf("beklenen 10, gerçek %d", stock.Available)
	}
}

func TestGetStock_Susccess(t *testing.T) {
	repo := &fakeStockRepo{
		stocks: make(map[int64]*domain.Stock),
	}

	svc := NewInventoryService(repo)

	_, err := svc.GetStock(context.Background(), 999)
	if err == nil {
		t.Fatal("stok bulunamadı hatası bekleniyordu")
	}

}

func TestReserveStock_Success(t *testing.T) {
	repo := &fakeStockRepo{stocks: make(map[int64]*domain.Stock)}
	svc := NewInventoryService(repo)

	svc.SetStock(context.Background(), 1, 10)
	err := svc.ReserveStock(context.Background(), 1, 3)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	stock, _ := svc.GetStock(context.Background(), 1)
	if stock.Available != 7 {
		t.Fatalf("beklenen 7, gerçek %d", stock.Available)
	}
}

func TestReserveStock_Insufficient(t *testing.T) {
	repo := &fakeStockRepo{stocks: make(map[int64]*domain.Stock)}
	svc := NewInventoryService(repo)

	svc.SetStock(context.Background(), 1, 3)
	err := svc.ReserveStock(context.Background(), 1, 5)

	if err == nil {
		t.Fatal("yetersiz stok hatası bekleniyordu")
	}

}

func TestReserveStock_NotFound(t *testing.T) {
	repo := &fakeStockRepo{stocks: map[int64]*domain.Stock{}}
	svc := NewInventoryService(repo)

	err := svc.ReserveStock(context.Background(), 999, 1)
	if err == nil {
		t.Fatal("stok bulunamadı hatası bekleniyordu")
	}

}

func TestReleaseStock(t *testing.T) {
	repo := &fakeStockRepo{stocks: map[int64]*domain.Stock{}}
	svc := NewInventoryService(repo)

	svc.SetStock(context.Background(), 1, 5)
	err := svc.ReleaseStock(context.Background(), 1, 3)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	stock, _ := svc.GetStock(context.Background(), 1)
	if stock.Available != 8 {
		t.Fatalf("beklenen 8, gerçek %d", stock.Available)
	}

}
