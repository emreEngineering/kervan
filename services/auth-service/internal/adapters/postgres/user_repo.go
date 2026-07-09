package postgres

import (
	"context"
	"database/sql"

	"github.com/emreEngineering/kervan/services/auth-service/internal/domain"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		db: db,
	}
}

func (r *PostgresUserRepo) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users(email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.CreatedAt).Scan(&id)

	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email=$1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		// kullanıcı bulunamadı hatası
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
