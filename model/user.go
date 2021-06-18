package model

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id"  db:"id"`
	Name      string    `json:"name"  db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserIdInfo struct {
	ID   string `json:"id"  db:"id"`
	Name string `json:"name"  db:"name"`
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
	_, err := db.ExecContext(ctx, "INSERT IGNORE INTO users (id, name) VALUES (?, ?) ", user.ID, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(ctx context.Context, accessToken string) ([]*UserIdInfo, error) {
	userId := []*UserIdInfo{}
	err := db.SelectContext(ctx, &userId, "SELECT * FROM users WHERE ID=c.Param(\"userID\") ")
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}

	return userId, nil
}
