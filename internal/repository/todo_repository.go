package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO todo_api (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func GetAllTodo(pool *pgxpool.Pool) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todo_api
		ORDER by created_at DESC;
	`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
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

func GetTodoById(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todo_api
		WHERE id = $1;
	`
	var todo models.Todo
	var err = pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, complete bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE todo_api
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, title, completed, updated_at, created_at;
	`

	var todo models.Todo

	var err = pool.QueryRow(ctx, query, title, complete, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.UpdatedAt,
		&todo.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func DeleteTodo(pool *pgxpool.Pool, id int) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE from todo_api
		WHERE id = $1
		
	`

	commandTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return nil
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo with id %d was not found", id)
	}

	return nil

}
