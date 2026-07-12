package domain

import "testing"

func TestNewProduct_Success(t *testing.T) {
	product, err := NewProduct("Test Ürün", "Açıklama", 99.99, "Elektronik")

	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}
	if product == nil {
		t.Fatal("product nil olamaz")
	}
	if product.Name != "Test Ürün" {
		t.Errorf("name = %s, beklenen Test Ürün", product.Name)
	}
	if product.Price != 99.99 {
		t.Errorf("price = %f, beklenen 99.99", product.Price)
	}
}

func TestNewProduct_EmptyName(t *testing.T) {
	_, err := NewProduct("", "Açıklama", 99.99, "Elektronik")
	if err == nil {
		t.Fatal("boş ürün adı için hata bekleniyordu")
	}
}

func TestNewProduct_InvalidPrice(t *testing.T) {
	_, err := NewProduct("Test Ürün", "Açıklama", 0.0, "Elektronik")
	if err == nil {
		t.Fatal("sıfır fiyat için hata bekleniyordu")
	}
}
