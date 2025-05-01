package model

import "time"

type Institution struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Type      string  `json:"type"`
	State     string  `json:"state"`
	City      *string `json:"city"`
	SourceUrl *string `json:"sourceUrl"`
	Active    bool    `json:"active"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
