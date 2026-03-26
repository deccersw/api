package service

import (
	"context"
	"errors"
	"todo_api/internal/domain"
	"todo_api/internal/ports"
)

type todoService struct {
	repo ports.TodoRepository
}

func NewTodoService(repo ports.TodoRepository) ports.TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(ctx context.Context, input domain.CreateTodoInput, userID string) (*domain.Todo, error) {
	if input.Title == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.Create(ctx, input, userID)
}

func (s *todoService) GetAll(ctx context.Context, userID string) ([]domain.Todo, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *todoService) GetByID(ctx context.Context, id int, userID string) (*domain.Todo, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *todoService) Update(ctx context.Context, id int, input domain.UpdateTodoInput, userID string) (*domain.Todo, error) {
	existing, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if input.Title == nil {
		input.Title = &existing.Title
	}
	if input.Completed == nil {
		input.Completed = &existing.Completed
	}

	return s.repo.Update(ctx, id, input, userID)
}

func (s *todoService) Delete(ctx context.Context, id int, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
