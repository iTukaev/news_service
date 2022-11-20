package models

import "time"

type News struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Article     string
}
