package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*domain.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO todo_api (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id;
	`

	var todo domain.Todo

	var err error = pool.QueryRow(ctx, query, title, completed, userID).Scan(
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

func GetAllTodo(pool *pgxpool.Pool, userID string) ([]domain.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todo_api
		WHERE user_id = $1
		ORDER by created_at DESC;
	`

	var rows, err = pool.Query(ctx, query, userID)

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

func GetTodoById(pool *pgxpool.Pool, id int, userID string) (*domain.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todo_api
		WHERE id = $1 AND user_id = $2;
	`
	var todo domain.Todo
	var err = pool.QueryRow(ctx, query, id, userID).Scan(
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

func UpdateTodo(pool *pgxpool.Pool, id int, title string, complete bool, userID string) (*domain.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE todo_api
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $3
		RETURNING id, title, completed, updated_at, created_at, user_id;
	`

	var todo domain.Todo

	var err = pool.QueryRow(ctx, query, title, complete, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UpdatedAt,
		&todo.CreatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE from todo_api
		WHERE id = $1 AND user_id = $2;
	`

	commandTag, err := pool.Exec(ctx, query, id, userID)
	if err != nil {
		return nil
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo with id %d was not found", id)
	}

	return nil

}
