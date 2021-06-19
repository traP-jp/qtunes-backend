package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Favorite struct {
	UserID     string    `json:"user_id"  db:"user_id"`
	ComposerID string    `json:"composer_id"  db:"composer_id"`
	SoundID    string    `json:"sound_id"  db:"sound_id"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}

type FavoriteCount struct {
	SoundID string `db:"sound_id"`
	Count   uint32 `db:"count"`
}

func getFavoriteCounts(ctx context.Context) (map[string]uint32, error) {
	var favCnt []*FavoriteCount
	err := db.SelectContext(ctx, &favCnt, "SELECT sound_id, COUNT( sound_id ) AS count FROM favorites GROUP BY sound_id")
	if err != nil {
		return nil, fmt.Errorf("Failed to get favorite counts: %w", err)
	}
	res := make(map[string]uint32)
	for _, v := range favCnt {
		res[v.SoundID] = v.Count
	}

	return res, nil
}

func getFavoriteCount(ctx context.Context, fileID string) (*FavoriteCount, error) {
	favCnt := FavoriteCount{}
	err := db.GetContext(ctx, &favCnt, "SELECT COUNT( composer_id ) AS count FROM favorites WHERE sound_id = ? LIMIT 1", fileID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get favorite count: %w", err)
	}

	return &favCnt, nil
}

func getMyFavorites(ctx context.Context, userID string) (map[string]bool, error) {
	var myFavorites []string
	err := db.SelectContext(ctx, &myFavorites, "SELECT sound_id FROM favorites WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Your Favorite Files: %w", err)
	}

	res := make(map[string]bool)
	for _, v := range myFavorites {
		res[v] = true
	}

	return res, nil
}

func getMyFavorite(ctx context.Context, userID string, fileID string) (bool, error) {
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

func insertFileFavorite(ctx context.Context, userID string, composerID string, fileID string) error {
	var _flag string
	err := db.GetContext(ctx, &_flag, "SELECT sound_id FROM favorites WHERE user_id = ? AND sound_id = ? LIMIT 1", userID, fileID)
	if err != nil && err != DBErrs["ErrNoRows"] {
		return fmt.Errorf("Failed to toggle favorite: %w", err)
	}
	if err == nil { // 既にfavoriteしている
		return DBErrs["NoChange"]
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
	if err == DBErrs["ErrNoRows"] { // 元々favoriteしていない
		return DBErrs["NoChange"]
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
