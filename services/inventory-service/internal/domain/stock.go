package domain

import "errors"

type Stock struct {
	ProductID int64
	Available int32
}

func NewStock(productID int64, available int32) (*Stock, error) {
	if available < 0 {
		return nil, errors.New("stok miktarı negatif olamaz")
	}
	return &Stock{
		ProductID: productID,
		Available: available,
	}, nil
}

// yeterli stok varsa düşür, yoksa hata dön
func (s *Stock) Reserve(quantity int32) error {
	if quantity <= 0 {
		return errors.New("rezerve miktarı pozitif olmalı")
	}
	if s.Available < quantity {
		return errors.New("yetersiz stok")
	}
	s.Available -= quantity
	return nil
}

// stoku geri ekle (her zaman başarılı, çünkü stok iadesi engellenemez)
func (s *Stock) Release(quantity int32) {
	s.Available += quantity
}
