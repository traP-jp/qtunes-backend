package domain

type Composer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
}
