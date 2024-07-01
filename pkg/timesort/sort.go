package timesort

import (
	service_models "faceit_parser/internal/models/service"
)

type ByTime []service_models.MatchStats

func (t ByTime) Len() int           { return len(t) }
func (t ByTime) Less(i, j int) bool { return t[i].CreatedAt.Before(t[j].CreatedAt) }
func (t ByTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
