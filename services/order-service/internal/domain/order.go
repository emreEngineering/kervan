package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusReturned  OrderStatus = "returned"
)

type Order struct {
	ID          int64
	UserID      int64
	Items       []OrderItem
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
}

type OrderItem struct {
	ProductID int64
	Quantity  int32
	Price     float64
}

func NewOrder(userID int64, items []OrderItem) *Order {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return &Order{
		UserID:      userID,
		Items:       items,
		TotalAmount: total,
		Status:      OrderStatusPending,
		CreatedAt:   time.Now(),
	}
}

func (o *Order) MarkPaid() {
	o.Status = OrderStatusPaid
}
func (o *Order) MarkCancelled() {
	o.Status = OrderStatusCancelled
}
