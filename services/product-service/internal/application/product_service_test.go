package application

import (
	"context"
	"testing"

	"github.com/emreEngineering/kervan/services/product-service/internal/domain"
)

type fakeProductRepo struct {
	products map[int64]*domain.Product
	nextID   int64
}

func newFakeProductRepo() *fakeProductRepo {
	return &fakeProductRepo{
		products: make(map[int64]*domain.Product),
		nextID:   1,
	}
}

func (r *fakeProductRepo) Create(ctx context.Context, product *domain.Product) error {
	product.ID = r.nextID
	r.products[r.nextID] = product
	r.nextID++
	return nil
}

func (r *fakeProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, nil
	}
	return product, nil
}

func (r *fakeProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	result := make([]*domain.Product, 0, len(r.products))
	for _, p := range r.products {
		result = append(result, p)
	}
	return result, nil
}

func (r *fakeProductRepo) Update(ctx context.Context, product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *fakeProductRepo) Delete(ctx context.Context, id int64) error {
	delete(r.products, id)
	return nil
}

func TestCreateProduct_Success(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	product, err := svc.CreateProduct(context.Background(), "Test Ürün", "Açıklama", 99.99, "Elektronik")

	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if product == nil {
		t.Fatal("product nil olamaz")
	}
	if product.ID != 1 {
		t.Errorf("ID=%d beklenen 1", product.ID)
	}
	if product.Name != "Test Ürün" {
		t.Errorf("Name=%s beklene Test Ürün", product.Name)
	}
}

func TestGetProduct_Success(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	created, _ := svc.CreateProduct(context.Background(), "Test", "Açıklama", 10.0, "Kat")
	product, err := svc.GetProduct(context.Background(), created.ID)

	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if product.ID != created.ID {
		t.Errorf("ID=%d, beklenen %d", product.ID, created.ID)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	_, err := svc.GetProduct(context.Background(), 999)

	if err == nil {
		t.Fatal("olamayan ürün için hata bekleniyordu")
	}
}

func TestCreateProduct_EmptyName(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	_, err := svc.CreateProduct(context.Background(), "", "Açıklama", 99.99, "Elektronik")

	if err == nil {
		t.Fatal("boş ürün için hata bekleniyordu")
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	created, _ := svc.CreateProduct(context.Background(), "Eski", "Eski Açıklama", 10.0, "Eski Kat")

	updated, err := svc.UpdateProduct(context.Background(), created.ID, "Yeni", "Yeni Açıklama", 20.0, "Yeni Kat")
	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if updated.Name != "Yeni" || updated.Price != 20.0 {
		t.Errorf("güncelleme başarısız: %+v", updated)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	_, err := svc.UpdateProduct(context.Background(), 999, "Yeni", "Açıklama", 20.0, "Kat")
	if err == nil {
		t.Fatal("olmayan ürün için hata bekleniyordu")
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	created, _ := svc.CreateProduct(context.Background(), "Silinecek", "Açıklama", 10.0, "Kat")

	err := svc.DeleteProduct(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}

	// Silindiğini doğrula
	_, err = svc.GetProduct(context.Background(), created.ID)
	if err == nil {
		t.Fatal("silinen ürün bulunmamalıydı")
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	err := svc.DeleteProduct(context.Background(), 999)
	if err == nil {
		t.Fatal("olmayan ürün için hata bekleniyordu")
	}
}

func TestListProducts(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	svc.CreateProduct(context.Background(), "A", "açıklama", 10.0, "Kat")
	svc.CreateProduct(context.Background(), "B", "açıklama", 20.0, "Kat")

	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("ürün sayısı = %d, beklenen 2", len(products))
	}
}
