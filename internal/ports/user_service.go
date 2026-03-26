package ports

import (
	"context"
	"todo_api/internal/domain"
)

type UserService interface {
	Register(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
}
