package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Favorite struct {
	UserID     string    `db:"user_id"`
	ComposerID string    `db:"composer_id"`
	SoundID    string    `db:"sound_id"`
	CreatedAt  time.Time `db:"created_at"`
	Count      int
}

func getFavoriteCounts(ctx context.Context) (map[string]int, error) {
	var favCount []*Favorite
	err := db.SelectContext(ctx, &favCount, "SELECT sound_id, COUNT( sound_id ) AS count FROM favorites GROUP BY sound_id")
	if err != nil {
		return nil, fmt.Errorf("Failed to get favorite counts: %w", err)
	}
	res := make(map[string]int)
	for _, v := range favCount {
		res[v.SoundID] = v.Count
	}

	return res, nil
}

func getFavoriteCount(ctx context.Context, fileID string) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT( composer_id ) AS count FROM favorites WHERE sound_id = ? LIMIT 1", fileID)
	if err != nil {
		return 0, fmt.Errorf("Failed to get favorite count: %w", err)
	}

	return count, nil
}

func getMyFavorites(ctx context.Context, userID string) ([]Favorite, error) {
	var myFavs []Favorite
	err := db.SelectContext(ctx, &myFavs, "SELECT sound_id, created_at FROM favorites WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err == sql.ErrNoRows {
		return []Favorite{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get Your Favorite Files: %w", err)
	}

	return myFavs, nil
}

func getMyFavoritesMap(ctx context.Context, userID string) (map[string]time.Time, error) {
	myFavs, err := getMyFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make(map[string]time.Time)
	for _, v := range myFavs {
		res[v.SoundID] = v.CreatedAt
	}

	return res, nil
}

func getMyFavorite(ctx context.Context, userID, fileID string) (bool, error) {
	myFavorite := ""
	err := db.GetContext(ctx, &myFavorite, "SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1", userID, fileID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Failed to get Your Favorite File: %w", err)
	}
	return (myFavorite != ""), nil
}

func insertFileFavorite(ctx context.Context, info Favorite) error {
	var check int
	err := db.SelectContext(ctx, &check, "SELECT EXISTS (SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1) AS check", info.UserID, info.SoundID)
	if err != nil {
		return fmt.Errorf("Failed to add the file to your favorite: %w", err)
	}
	if check == 1 { // 既にお気に入りされているとき
		return ErrNoChange
	}

	_, err = db.ExecContext(ctx, "INSERT INTO favorites (user_id, composer_id, sound_id) VALUES (?, ?, ?)", info.UserID, info.ComposerID, info.SoundID)
	if err != nil {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}

	return nil
}

func deleteFileFavorite(ctx context.Context, info Favorite) error {
	var check int
	err := db.SelectContext(ctx, &check, "SELECT EXISTS (SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1) AS check", info.UserID, info.SoundID)
	if err != nil {
		return fmt.Errorf("Failed to remove the file from your favorite: %w", err)
	}
	if check == 1 { // 元からお気に入りされていないとき
		return ErrNoChange
	}

	_, err = db.ExecContext(ctx, "DELETE FROM favorites WHERE user_id = ? AND sound_id = ?", info.UserID, info.SoundID)
	if err != nil {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}

	return nil
}
