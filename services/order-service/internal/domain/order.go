package domain

import "time"

type Order struct {
	ID          int64
	UserID      int64
	Items       []OrderItem
	TotalAmount float64
	Status      string
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
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
}

func (o *Order) MarkPaid() {
	o.Status = "paid"
}
func (o *Order) MarkCancelled() {
	o.Status = "cancelled"
}
