package postgres

import (
	"context"
	"database/sql"

	"github.com/emreEngineering/kervan/services/inventory-service/internal/domain"
)

type PostgresStockRepo struct {
	db *sql.DB
}

func NewPostgresStockRepo(db *sql.DB) *PostgresStockRepo {
	return &PostgresStockRepo{db: db}
}

func (r *PostgresStockRepo) GetStock(ctx context.Context, productID int64) (*domain.Stock, error) {
	row := r.db.QueryRowContext(ctx, "SELECT product_id,available FROM stocks WHERE product_id=$1", productID)
	stock := &domain.Stock{}

	err := row.Scan(&stock.ProductID, &stock.Available)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return stock, nil

}

func (r *PostgresStockRepo) Save(ctx context.Context, stock *domain.Stock) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO stocks (product_id, available) VALUES ($1, $2)", stock.ProductID, stock.Available)
	return err
}

func (r *PostgresStockRepo) Update(ctx context.Context, stock *domain.Stock) error {
	_, err := r.db.ExecContext(ctx, "UPDATE stocks SET available=$1 WHERE product_id=$2", stock.Available, stock.ProductID)
	return err
}
