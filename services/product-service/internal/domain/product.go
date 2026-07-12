package domain

import (
	"errors"
	"time"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Category    string
	CreatedAt   time.Time
}

func NewProduct(name, description string, price float64, category string) (*Product, error) {
	if name == "" {
		return nil, errors.New("ürün adı boş olamaz")
	}
	if price <= 0 {
		return nil, errors.New("ürün fiyatı sıfır veya sıfırdan küçük olamaz")
	}
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		CreatedAt:   time.Now(),
	}, nil
}
