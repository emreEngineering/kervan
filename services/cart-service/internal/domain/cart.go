package domain

import "errors"

type CartItem struct {
	ProductID int64
	Quantity  int32
}

type Cart struct {
	UserID int64
	Items  []CartItem
}

func NewCart(userID int64) *Cart {
	return &Cart{UserID: userID}
}

func (c *Cart) AddItem(productID int64, quantity int32) error {
	if quantity <= 0 {
		return errors.New("miktar pozitif olmalı")
	}

	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += quantity
			return nil
		}
	}
	c.Items = append(c.Items, CartItem{ProductID: productID, Quantity: quantity})
	return nil
}

func (c *Cart) RemoveItem(productID int64) error {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}
	return errors.New("ürün sepette bulunamadı")
}

func (c *Cart) Clear() {
	c.Items = []CartItem{}
}
