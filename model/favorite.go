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
