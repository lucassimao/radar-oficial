package model

import "time"

// UserSession represents a user's WhatsApp session
type UserSession struct {
	ID            int        `json:"id"`
	PhoneNumber   string     `json:"phone_number"`
	CreatedAt     time.Time  `json:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at"`
	State         *string    `json:"state"`
}