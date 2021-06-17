package domain

import "time"

type File struct {
	ID             string    `json:"id"`
	ComposerID     string    `json:"composerId"`
	FavoriteCount    int       `json:"favorite_count"`
	IsFavoriteByMe bool      `json:"isFavoriteByMe"`
	CreatedAt      time.Time `json:"created_at"`
}
