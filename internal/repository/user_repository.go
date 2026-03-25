package repository

import (
	"context"
	"time"
	"todo_api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *domain.User) (*domain.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO user_api (email,password)
		VALUES($1,$2)
		RETURNING id, email, created_at, updated_at;
	`

	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*domain.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
		id, email, password, created_at, updated_at
		FROM user_api
		WHERE email = $1;
	`

	var user domain.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*domain.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
		id, email,password, created_at, updated_at
		FROM user_api
		WHERE id = $1;
	`

	var user domain.User

	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
