package application

import (
	"context"
	"errors"

	"github.com/emreEngineering/kervan/services/product-service/internal/domain"
	"github.com/emreEngineering/kervan/services/product-service/internal/ports"
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description string, price float64, category string) (*domain.Product, error) {
	product, err := domain.NewProduct(name, description, price, category)
	if err != nil {
		return nil, err
	}
	err = s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("ürün bulunamadı")
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, name, description string, price float64, category string) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("ürün bulunamadı")
	}
	product.Name = name
	product.Description = description
	product.Price = price
	product.Category = category

	err = s.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New("ürün bulunamadı")
	}
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
