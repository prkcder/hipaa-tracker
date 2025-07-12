package models


import (
	"time"

)

type Event struct {
	ID int `json:"id"`
	EventType string `json:"event_type"`
	Payload map[string]any `json:"payload"`
	Sanitized bool `json:"sanitized"`
	CreatedAt time.Time `json:"create_at"`
}

