package model

import (
	"context"
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

func getMyFavorite(ctx context.Context, userID string, fileID string) (bool, error) {
	myFavorite := ""
	err := db.GetContext(ctx, &myFavorite, "SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1", userID, fileID)
	if err == ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Failed to get Your Favorite File: %w", err)
	}
	return (myFavorite != ""), nil
}

func insertFileFavorite(ctx context.Context, userID string, composerID string, fileID string) error {
	var _flag string
	err := db.GetContext(ctx, &_flag, "SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1", userID, fileID)
	if err != nil && err != ErrNoRows {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}
	if err == nil { // 既にfavoriteしている
		return ErrNoChange
	}

	_, err = db.ExecContext(ctx, "INSERT INTO favorites (user_id, composer_id, sound_id) VALUES (?, ?, ?)", userID, composerID, fileID)
	if err != nil {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}

	return nil
}

func deleteFileFavorite(ctx context.Context, userID, fileID string) error {
	var _flag string
	err := db.GetContext(ctx, &_flag, "SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1", userID, fileID)
	if err == ErrNoRows { // 元々favoriteしていない
		return ErrNoChange
	}
	if err != nil {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM favorites WHERE user_id = ? AND sound_id = ?", userID, fileID)
	if err != nil {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}

	return nil
}
