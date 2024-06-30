package models

import (
	"time"
)

type Player struct {
	ID        string                 `json:"id"`
	Nickname  string                 `json:"nickname"`
	CreatedAt time.Time              `json:"created_at"`
	Payload   map[string]interface{} `json:"payload"`
}
