package model

import (
	"context"
	"fmt"
)

type User struct {
	ID     string `json:"id"  db:"id"`
	Name   string `json:"name"  db:"name"`
	Status int8   `json:"status"  db:"status"`
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
	return nil //TODO
}
