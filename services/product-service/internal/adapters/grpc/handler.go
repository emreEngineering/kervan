package grpc

import (
	"context"

	productv1 "github.com/emreEngineering/kervan/gen/go/product/v1"
	"github.com/emreEngineering/kervan/services/product-service/internal/application"
)

type ProductHandler struct {
	productv1.UnimplementedProductServiceServer
	app *application.ProductService
}

func NewProductHandler(app *application.ProductService) *ProductHandler {
	return &ProductHandler{app: app}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	product, err := h.app.CreateProduct(ctx, req.Name, req.Description, req.Price, req.Category)
	if err != nil {
		return nil, err
	}
	return &productv1.CreateProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	product, err := h.app.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &productv1.GetProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	products, err := h.app.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	var pbProducts []*productv1.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &productv1.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
		})
	}

	return &productv1.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductResponse, error) {
	product, err := h.app.UpdateProduct(ctx, req.Id, req.Name, req.Description, req.Price, req.Category)
	if err != nil {
		return nil, err
	}
	return &productv1.UpdateProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*productv1.DeleteProductResponse, error) {
	err := h.app.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &productv1.DeleteProductResponse{
		Success: true,
	}, nil
}
