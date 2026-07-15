package domain

import "testing"

func TestNewCart(t *testing.T) {
	c := NewCart(1)
	if c.UserID != 1 {
		t.Fatalf("beklenen 1, gerçek %d", c.UserID)
	}

	if len(c.Items) != 0 {
		t.Fatalf("sepet boş olmalı, %d eleman var", len(c.Items))
	}

}

func TestAddItem_NewProduct(t *testing.T) {
	c := NewCart(1)
	err := c.AddItem(100, 2)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	if len(c.Items) != 1 || c.Items[0].Quantity != 2 {
		t.Fatalf("beklenen 1 ürün iki adet, gerçek %d ürün %d adet", len(c.Items), c.Items[0].Quantity)
	}
}

func TestAddItem_ExistingProduct(t *testing.T) {
	c := NewCart(1)
	c.AddItem(100, 2)
	c.AddItem(100, 3)

	if len(c.Items) != 1 || c.Items[0].Quantity != 5 {
		t.Fatalf("beklenen 1 ürün 5 adet, gerçek %d ürün %d adet", len(c.Items), c.Items[0].Quantity)
	}
}

func TestAddItem_InvalidQuantity(t *testing.T) {
	c := NewCart(1)

	err := c.AddItem(100, 0)
	if err == nil {
		t.Fatal("pozitif miktar hatası bekleniyordu")
	}
}

func TestRemoveItem_Success(t *testing.T) {
	c := NewCart(1)
	c.AddItem(100, 2)
	c.AddItem(200, 1)
	err := c.RemoveItem(100)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	if len(c.Items) != 1 || c.Items[0].ProductID != 200 {
		t.Fatalf("beklenen 1 ürün (200), gerçek %d ürün", len(c.Items))
	}
}

func TestRemoveItem_NotFound(t *testing.T) {
	c := NewCart(1)
	err := c.RemoveItem(999)
	if err == nil {
		t.Fatal("ürün bulunamadı hatası bekleniyordu")
	}
}

func TestClear(t *testing.T) {
	c := NewCart(1)
	c.AddItem(100, 2)
	c.AddItem(200, 1)
	c.Clear()
	if len(c.Items) != 0 {
		t.Fatalf("sepet boş olmalı, %d eleman var", len(c.Items))
	}
}
