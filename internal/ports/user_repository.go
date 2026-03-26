package ports

import (
	"context"
	"todo_api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
