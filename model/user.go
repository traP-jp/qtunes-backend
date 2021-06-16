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
	TraqID    string    `json:"traq_id"  db:"traq_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
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
	_, err := db.Query("INSERT INTO users (id, traq_id) VALUES (?, ?)", user.ID, user.TraqID)
	if err != nil {
		return err
	}

	return nil
}
