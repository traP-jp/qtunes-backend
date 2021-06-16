package model

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

var baseURL, _ = url.Parse("https://q.trap.jp/api/v3")

type User struct {
	ID        string    `json:"id"  db:"id"`
	Name      string    `json:"name"  db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}
	err := db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}

	return users, nil
}

func CreateUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = name", user.ID, user.Name)
	if err != nil {
		return err
	}

	return nil
}
