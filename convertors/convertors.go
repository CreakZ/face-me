package convertors

import (
	"faceit_parser/models"
	"strconv"
	"strings"
	"time"
)

func ConvertMatchStats(match models.MatchStatsRaw) models.MatchStats {
	kills, _ := strconv.ParseUint(match.Kills, 10, 8)
	assists, _ := strconv.ParseUint(match.Assists, 10, 8)
	deaths, _ := strconv.ParseUint(match.Deaths, 10, 8)
	headshots, _ := strconv.ParseUint(match.Headshots, 10, 8)

	score := strings.Split(match.Score, " / ")

	left, _ := strconv.ParseUint(score[0], 10, 8)
	right, _ := strconv.ParseUint(score[1], 10, 8)

	return models.MatchStats{
		CreatedAt: time.UnixMilli(match.CreatedAt),
		Score:     [2]uint8{uint8(left), uint8(right)},
		Map:       match.Map,
		Kills:     uint8(kills),
		Assists:   uint8(assists),
		Deaths:    uint8(deaths),
		Headshots: uint8(headshots),
	}
}
