package model

import (
	"context"
	"fmt"
)

type FavoriteCount struct {
	ComposerID string `db:"composer_id"`
	Count      uint8  `db:"count"`
}
type MyFavorites []string

func getFavoriteCounts(ctx context.Context) (map[string]uint8, error) {
	favCnt := make([]*FavoriteCount, 0)
	err := db.SelectContext(ctx, &favCnt, "SELECT composer_id, COUNT( composer_id ) as count FROM favorites GROUP BY composer_id")
	if err != nil {
		return nil, fmt.Errorf("Failed to get favorite counts: %w", err)
	}
	res := make(map[string]uint8)
	for _, v := range favCnt {
		res[v.ComposerID] = v.Count
	}

	return res, nil
}

func getMyFavorites(ctx context.Context, userID string) (map[string]bool, error) {
	myFavorites := make(MyFavorites, 0)
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
