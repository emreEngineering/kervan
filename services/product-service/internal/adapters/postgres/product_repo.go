package postgres

import (
	"context"
	"database/sql"

	"github.com/emreEngineering/kervan/services/product-service/internal/domain"
)

type PostgresProductRepo struct {
	db *sql.DB
}

func NewPostgresProductRepo(db *sql.DB) *PostgresProductRepo {
	return &PostgresProductRepo{db: db}
}

func (r *PostgresProductRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products(name,description,price,category,created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.Category, product.CreatedAt).Scan(&id)
	if err != nil {
		return err
	}
	product.ID = id
	return nil
}

func (r *PostgresProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, name, description, price, category, created_at FROM products WHERE id=$1`
	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Category,
		&product.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *PostgresProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, name, description, price, category, created_at FROM products ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *PostgresProductRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, category = $4 WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Category, product.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresProductRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
