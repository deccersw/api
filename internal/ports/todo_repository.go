package ports

import (
	"context"
	"todo_api/internal/domain"
)

type TodoRepository interface {
	Create(ctx context.Context, input domain.CreateTodoInput, userID string) (*domain.Todo, error)
	GetAll(ctx context.Context, userID string) ([]domain.Todo, error)
	GetByID(ctx context.Context, id int, userID string) (*domain.Todo, error)
	Update(ctx context.Context, id int, input domain.UpdateTodoInput, userID string) (*domain.Todo, error)
	Delete(ctx context.Context, id int, userID string) error
}
