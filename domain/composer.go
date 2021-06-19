package domain

import "time"

type Composer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
}
type ComposerFile struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	ComposerID     *string   `json:"composer_id"`
	ComposerName   string    `json:"composer_name"`
	FavoriteCount  uint32    `json:"favorite_count"`
	IsFavoriteByMe bool      `json:"is_favorite_by_me"`
	CreatedAt      time.Time `json:"created_at"`
}
