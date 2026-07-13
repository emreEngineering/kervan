package grpc

import (
	"context"

	inventoryv1 "github.com/emreEngineering/kervan/gen/go/inventory/v1"
	"github.com/emreEngineering/kervan/services/inventory-service/internal/application"
)

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	app *application.InventoryService
}

func NewInventoryHandler(app *application.InventoryService) *InventoryHandler {
	return &InventoryHandler{app: app}
}

func (h *InventoryHandler) GetStock(ctx context.Context, req *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error) {
	stock, err := h.app.GetStock(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &inventoryv1.GetStockResponse{
		Available: stock.Available,
	}, nil
}

func (h *InventoryHandler) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	err := h.app.ReserveStock(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &inventoryv1.ReserveStockResponse{Success: true}, nil
}

func (h *InventoryHandler) ReleaseStock(ctx context.Context, req *inventoryv1.ReleaseStockRequest) (*inventoryv1.ReleaseStockResponse, error) {
	err := h.app.ReleaseStock(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ReleaseStockResponse{Success: true}, nil
}

func (h *InventoryHandler) SetStock(ctx context.Context, req *inventoryv1.SetStockRequest) (*inventoryv1.SetStockResponse, error) {
	_, err := h.app.SetStock(ctx, req.ProductId, req.Available)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.SetStockResponse{Success: true}, nil
}
