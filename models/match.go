package models

import "time"

// Struct that represents stats of 1 match
type MatchStats struct {
	CreatedAt time.Time `json:"created_at"`
	Score     [2]uint8  `json:"i18"`
	Map       string    `json:"i1"`
	Kills     uint8     `json:"i6"`
	Assists   uint8     `json:"i7"`
	Deaths    uint8     `json:"i8"`
	Headshots uint8     `json:"i13"`
}
