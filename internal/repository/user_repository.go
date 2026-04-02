package repository

import (
	"context"
	"errors"
	"todo_api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {

	query := `
		INSERT INTO user_api (email,password)
		VALUES($1,$2)
		RETURNING id, email, created_at, updated_at;
	`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, input.Email, input.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	query := `
		SELECT 
		id, email, password, created_at, updated_at
		FROM user_api
		WHERE email = $1;
	`

	var user domain.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {

	query := `
		SELECT 
		id, email,password, created_at, updated_at
		FROM user_api
		WHERE id = $1;
	`

	var user domain.User

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
