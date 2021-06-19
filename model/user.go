package model

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id"  db:"id"`
	Name      string    `json:"name"  db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type UserMe struct {
	ID            string   `json:"id"  db:"id"`
	Name          string   `json:"name"  db:"name"`
	FavoriteFiles []string `json:"favoriteFiles" db:"sound_id"`
}

func GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}
	err := db.SelectContext(ctx, &users, "SELECT id, name FROM users")
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

func GetUser(ctx context.Context, accessToken string, userID string) (*User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, name FROM users WHERE id = ? LIMIT 1", userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	return &user, nil
}

func GetUserMe(ctx context.Context, accessToken string) (*UserMe, error) {
	var userMe UserMe
	err := db.GetContext(ctx, &userMe, "SELECT id, name, sound_id FROM users FULL JOIN favorites  LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("Failed to get userMe: %w", err)
	}

	return &userMe, nil
}
