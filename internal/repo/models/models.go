package models

import "time"

type (
	ListParams struct {
		Limit  uint64 `json:"limit"`
		Offset uint64 `json:"offset"`
		Order  bool   `json:"order"`
	}

	News struct {
		ID          uint32
		Title       string
		Link        string
		Description string
		PubDate     time.Time
		Article     string
	}
)
