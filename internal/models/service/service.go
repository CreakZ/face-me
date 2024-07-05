package service_models

import "time"

// Struct that represents stats of 1 match
type MatchStats struct {
	CreatedAt time.Time `json:"created_at"`
	Score     [2]uint8  `json:"score"`
	Map       string    `json:"map"`
	Kills     uint8     `json:"kills"`
	Assists   uint8     `json:"assists"`
	Deaths    uint8     `json:"deaths"`
	Headshots uint8     `json:"headshots"`
}
