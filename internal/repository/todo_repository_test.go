package repository

import (
	"context"
	"os"
	"testing"

	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dbURL := "postgres://postgres:password@localhost:5432/todo_api"

	var err error
	testPool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	testPool.Close()
	os.Exit(code)
}

func cleanDB(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(context.Background(),
		"DELETE FROM todos")
	if err != nil {
		t.Fatal(err)
	}
}

func insertTodo(t *testing.T) models.Todo {
	t.Helper()

	var todo models.Todo

	err := testPool.QueryRow(
		context.Background(),
		`INSERT INTO todos(title, completed, user_id)
		VALUES($1,$2,$3)
		RETURNING id,title,completed,created_at,updated_at,user_id`,
		"Test Todo",
		false,
		"user1",
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		t.Fatal(err)
	}

	return todo
}

func TestCreateTodo(t *testing.T) {
	cleanDB(t)

	todo, err := CreateTodo(testPool, "Learn Go", false, "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.Title != "Learn Go" {
		t.Errorf("expected Learn Go, got %s", todo.Title)
	}

	if todo.UserID != "user1" {
		t.Errorf("expected user1 got %s", todo.UserID)
	}
}

func TestGetTodoByID(t *testing.T) {
	cleanDB(t)

	inserted := insertTodo(t)

	todo, err := GetToDoByID(testPool, inserted.ID, "user1")

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if todo.ID != inserted.ID {
		t.Errorf("expected %d got %d", inserted.ID, todo.ID)
	}
}

func TestGetAllTodos(t *testing.T) {
	cleanDB(t)

	insertTodo(t)
	insertTodo(t)

	todos, err := GetAllTodos(testPool, "user1")

	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 2 {
		t.Errorf("expected 2 todos got %d", len(todos))
	}
}

func TestUpdateTodo(t *testing.T) {
	cleanDB(t)

	inserted := insertTodo(t)

	updated, err := UpdateToDo(
		testPool,
		inserted.ID,
		"Updated Todo",
		true,
		"user1",
	)

	if err != nil {
		t.Fatal(err)
	}

	if updated.Title != "Updated Todo" {
		t.Errorf("expected Updated Todo got %s", updated.Title)
	}

	if !updated.Completed {
		t.Error("expected completed=true")
	}
}

func TestDeleteTodo(t *testing.T) {
	cleanDB(t)

	inserted := insertTodo(t)

	err := DeleteToDo(testPool, inserted.ID, "user1")

	if err != nil {
		t.Fatal(err)
	}

	_, err = GetToDoByID(testPool, inserted.ID, "user1")

	if err == nil {
		t.Error("expected todo to be deleted")
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	cleanDB(t)

	err := DeleteToDo(testPool, 9999, "user1")

	if err == nil {
		t.Error("expected error")
	}
}
