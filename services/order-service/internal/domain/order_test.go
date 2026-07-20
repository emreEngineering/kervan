package domain

import (
	"testing"
)

func TestNewOrder_Success(t *testing.T) {
	items := []OrderItem{
		{ProductID: 100, Quantity: 2, Price: 50},
		{ProductID: 200, Quantity: 1, Price: 100},
	}

	order := NewOrder(1, items)
	if order.UserID != 1 {
		t.Fatalf("beklenen UserID=1, gerçek= %d", order.ID)
	}
	if order.TotalAmount != 200 {
		t.Fatalf("beklenen TotalAmount=200, gerçek= %f", order.TotalAmount)
	}
	if order.Status != OrderStatusPending {
		t.Fatalf("beklenen status=pending, gerçek= %s", order.Status)
	}
	if len(order.Items) != 2 {
		t.Fatalf("beklenen 2 ürün, gerçek= %d", len(order.Items))
	}
}

func TestNewOrder_EmptyItems(t *testing.T) {
	order := NewOrder(1, []OrderItem{})
	if order.TotalAmount != 0 {
		t.Fatalf("beklenen TotalAmount=0, gerçek= %f", order.TotalAmount)
	}
	if len(order.Items) != 0 {
		t.Fatalf("beklenen 0 ürün, gerçek= %d", len(order.Items))
	}
}

func TestMarkPaid(t *testing.T) {
	items := []OrderItem{{ProductID: 100, Quantity: 1, Price: 50}}
	order := NewOrder(1, items)
	order.MarkPaid()
	if order.Status != OrderStatusPaid {
		t.Fatalf("beklenen paid, gerçek= %s", order.Status)
	}
}

func TestMarkCancelled(t *testing.T) {
	items := []OrderItem{{ProductID: 100, Quantity: 1, Price: 50}}
	order := NewOrder(1, items)
	order.MarkCancelled()
	if order.Status != OrderStatusCancelled {
		t.Fatalf("beklenen cancelled, gerçek= %s", order.Status)
	}
}
