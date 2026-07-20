package postgres

import (
	"context"
	"database/sql"

	"github.com/emreEngineering/kervan/services/order-service/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Save(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "INSERT INTO orders (user_id, total_amount, status) VALUES ($1,$2,$3) RETURNING id", order.UserID, order.TotalAmount, order.Status).Scan(&order.ID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, "INSERT INTO order_items (order_id,product_id,quantity,price) VALUES ($1,$2,$3,$4)", order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}

	}
	return tx.Commit()
}

func (r *OrderRepo) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	order := &domain.Order{}
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, total_amount, status, created_at FROM orders WHERE id=$1", id).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, "SELECT product_id, quantity, price FROM order_items WHERE order_id=$1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	// döngü sırasında bağlantı koparsa veri yutulmasını önlemek için
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return order, nil

}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus) error {
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status=$1 WHERE id=$2", status, id)
	return err
}
