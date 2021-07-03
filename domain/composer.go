package domain

import "time"

type Composer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
	UpdatedAt time.Time
}
