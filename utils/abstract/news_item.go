package abstract

import (
	"time"
)

type NewsItem struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Url       string    `json:"url"`
	Author    string    `json:"author"`
	Category  string    `json:"category"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ByCreatedAt []NewsItem

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }
