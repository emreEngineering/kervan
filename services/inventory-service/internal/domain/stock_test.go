package domain

import "testing"

func TestNewStock_Success(t *testing.T) {
	s, err := NewStock(1, 10)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if s.ProductID != 1 || s.Available != 10 {
		t.Fatalf("yanlış değerler: %+v", s)
	}
}

func TestNewStock_NegativeQuantity(t *testing.T) {
	_, err := NewStock(1, -5)
	if err == nil {
		t.Fatal("negatif stok hatası bekleniyordu")
	}
}

func TestReserve_Success(t *testing.T) {
	s, _ := NewStock(1, 10)
	err := s.Reserve(3)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if s.Available != 7 {
		t.Fatalf("beklenen 7, gerçek %d", s.Available)
	}
}

func TestReserve_InsufficientStock(t *testing.T) {
	s, _ := NewStock(1, 3)
	err := s.Reserve(5)
	if err == nil {
		t.Fatal("yetersiz stok hatası bekleniyordu")
	}
}

func TestReserve_NegativeQuantity(t *testing.T) {
	s, _ := NewStock(1, 10)
	err := s.Reserve(-1)
	if err == nil {
		t.Fatal("negatif miktar hatası bekleniyordu")
	}
}

func TestRelease(t *testing.T) {
	s, _ := NewStock(1, 5)
	s.Release(3)
	if s.Available != 8 {
		t.Fatalf("beklenen 8, gerçek %d", s.Available)
	}
}
