package model

import (
	"context"
	"time"
)

type Favorite struct { //TODO db
	UserID string
	ComposerID string
	SoundID string
	CreatedAt time.Time
}

type FavoriteCounts map[string]int

func getFavoriteCounts(ctx context.Context) (FavoriteCount, error){
	favoriteCounts := make(FavoriteCounts)
	db.ExecContext(ctx, &favoriteCounts, "SELECT COUNT (composer_id) FROM users")
	return nil, nil //TODO: 角
}

