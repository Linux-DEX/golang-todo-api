package repository

import (
	"context"
	"testing"

	"todo_api/internal/models"
)

func cleanUsersTable(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	cleanUsersTable(t)

	user := &models.User{
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	createdUser, err := CreateUser(testPool, user)
	if err != nil {
		t.Fatalf("CreateUser() returned error: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be set")
	}

	if createdUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, createdUser.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	cleanUsersTable(t)

	user := &models.User{
		Email:    "alice@example.com",
		Password: "hashedpassword",
	}

	_, err := CreateUser(testPool, user)
	if err != nil {
		t.Fatal(err)
	}

	foundUser, err := GetUserByEmail(testPool, user.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail() returned error: %v", err)
	}

	if foundUser.Email != user.Email {
		t.Errorf("expected %s, got %s", user.Email, foundUser.Email)
	}
}

func TestGetUserByID(t *testing.T) {
	cleanUsersTable(t)

	user := &models.User{
		Email:    "bob@example.com",
		Password: "hashedpassword",
	}

	createdUser, err := CreateUser(testPool, user)
	if err != nil {
		t.Fatal(err)
	}

	foundUser, err := GetUserByID(testPool, createdUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID() returned error: %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("expected ID %s, got %s", createdUser.ID, foundUser.ID)
	}

	if foundUser.Email != createdUser.Email {
		t.Errorf("expected email %s, got %s", createdUser.Email, foundUser.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	cleanUsersTable(t)

	_, err := GetUserByEmail(testPool, "notfound@example.com")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	cleanUsersTable(t)

	_, err := GetUserByID(testPool, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}
