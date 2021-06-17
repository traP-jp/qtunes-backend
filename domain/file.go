package domain

import "time"

type File struct {
	ID             string    `json:"id"`
	ComposerID     string    `json:"composer_id"`
	FavoriteCount  uint32     `json:"favorite_count"`
	IsFavoriteByMe bool      `json:"is_favorite_by_me"`
	CreatedAt      time.Time `json:"created_at"`
}
