package repository

import (
	"context"
	"errors"
	"todo_api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoRepository(pool *pgxpool.Pool) *todoRepository {
	return &todoRepository{pool: pool}
}

func (r *todoRepository) Create(ctx context.Context, input domain.CreateTodoInput, userID string) (*domain.Todo, error) {

	var query string = `
		INSERT INTO todo_api (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id;
	`

	var todo domain.Todo

	var err error = r.pool.QueryRow(ctx, query, input.Title, input.Completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func (r *todoRepository) GetAll(ctx context.Context, userID string) ([]domain.Todo, error) {

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todo_api
		WHERE user_id = $1
		ORDER by created_at DESC;
	`

	var rows, err = r.pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []domain.Todo = []domain.Todo{}

	for rows.Next() {
		var todo domain.Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil

}

func (r *todoRepository) GetByID(ctx context.Context, id int, userID string) (*domain.Todo, error) {

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todo_api
		WHERE id = $1 AND user_id = $2;
	`
	var todo domain.Todo
	var err = r.pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &todo, nil

}

func (r *todoRepository) Update(ctx context.Context, id int, input domain.UpdateTodoInput, userID string) (*domain.Todo, error) {

	var query string = `
		UPDATE todo_api
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, completed, updated_at, created_at, user_id;
	`

	var todo domain.Todo

	var err = r.pool.QueryRow(ctx, query, input.Title, input.Completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UpdatedAt,
		&todo.CreatedAt,
		&todo.UserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &todo, nil

}

func (r *todoRepository) Delete(ctx context.Context, id int, userID string) error {

	var query string = `
		DELETE from todo_api
		WHERE id = $1 AND user_id = $2;
	`

	commandTag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil

}
