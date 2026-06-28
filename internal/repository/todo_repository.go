package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllTodos(pool *pgxpool.Pool, userID string) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id, title, completed, created_at, updated_at, user_id
	FROM todos
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []models.Todo{}

	for rows.Next() {
		todo := models.Todo{}

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

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO todos (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo

	err := pool.QueryRow(ctx, query, title, completed, userID).Scan(
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

func GetToDoByID(pool *pgxpool.Pool, id int, userID string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	todo := models.Todo{}

	err := pool.QueryRow(ctx, query, id, userID).Scan(
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

func UpdateToDo(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE todos
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	todo := models.Todo{}

	err := pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
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

func DeleteToDo(pool *pgxpool.Pool, id int, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM todos
		WHERE id = $1 AND user_id = $2
	`

	commandTag, err := pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}
