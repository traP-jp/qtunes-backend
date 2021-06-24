package domain

import "time"

type File struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	ComposerID     string    `json:"composer_id"`
	ComposerName   string    `json:"composer_name"`
	MessageID      string    `json:"message_id"`
	FavoriteCount  int       `json:"favorite_count"`
	IsFavoriteByMe bool      `json:"is_favorite_by_me"`
	CreatedAt      time.Time `json:"created_at"`
}
