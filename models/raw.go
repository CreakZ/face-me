package models

type MatchStatsRaw struct {
	CreatedAt int64  `json:"created_at"`
	Score     string `json:"i18"`
	Map       string `json:"i1"`
	Kills     string `json:"i6"`
	Assists   string `json:"i7"`
	Deaths    string `json:"i8"`
	Headshots string `json:"i13"`
}
