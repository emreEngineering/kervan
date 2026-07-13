package application

import (
	"context"
	"errors"

	"github.com/emreEngineering/kervan/services/inventory-service/internal/domain"
	"github.com/emreEngineering/kervan/services/inventory-service/internal/ports"
)

type InventoryService struct {
	repo ports.StockRepository
}

func NewInventoryService(repo ports.StockRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetStock(ctx context.Context, productID int64) (*domain.Stock, error) {
	stock, err := s.repo.GetStock(ctx, productID)

	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, errors.New("stok bulunamadı")
	}
	return stock, nil
}

func (s *InventoryService) ReserveStock(ctx context.Context, productID int64, quantity int32) error {
	stock, err := s.repo.GetStock(ctx, productID)
	if err != nil {
		return err
	}
	if stock == nil {
		return errors.New("stok bulunamadı")
	}
	if err := stock.Reserve(quantity); err != nil {
		return err
	}
	return s.repo.Update(ctx, stock)
}

func (s *InventoryService) ReleaseStock(ctx context.Context, productID int64, quantity int32) error {
	stock, err := s.repo.GetStock(ctx, productID)
	if err != nil {
		return err
	}
	if stock == nil {
		return errors.New("stok bulunamadı")
	}
	stock.Release(quantity)
	return s.repo.Update(ctx, stock)
}

func (s *InventoryService) SetStock(ctx context.Context, productID int64, available int32) (*domain.Stock, error) {
	stock, err := domain.NewStock(productID, available)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, stock); err != nil {
		return nil, err
	}
	return stock, nil
}
