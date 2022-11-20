package client

import "time"

type News struct {
	Title       string    `json:"title,omitempty"`
	Link        string    `json:"link,omitempty"`
	Description string    `json:"description,omitempty"`
	PubDate     time.Time `json:"pubDate"`
	Article     string    `json:"article,omitempty"`
}

type Page struct {
	Size uint64 `json:"size,omitempty"`
	Page uint64 `json:"page,omitempty"`
}

type NewsList struct {
	List []News `json:"list"`
	Page Page   `json:"page"`
}
