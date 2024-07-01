package service

import (
	"context"
	"faceit_parser/internal/convertors"
	service_models "faceit_parser/internal/models/service"
	"faceit_parser/internal/repository"
	"faceit_parser/pkg/timesort"
	"log"
	"sort"
)

type service struct {
	repo   repository.Repository
	logger *log.Logger
}

func InitService(repo repository.Repository) service {
	return service{
		repo:   repo,
		logger: log.Default(),
	}
}

type Service interface {
	GetAllMatchesStats(context context.Context, nickname, game string) ([]service_models.MatchStats, error)
	GetMatchWithMostKills(ctx context.Context, matches []service_models.MatchStats) (service_models.MatchStats, error)
	GetLongestMatch(ctx context.Context, matches []service_models.MatchStats) (service_models.MatchStats, error)
	GetGlobalKD(ctx context.Context, matches []service_models.MatchStats) (float64, error)
}

func (s service) GetAllMatchesStats(context context.Context, nickname, game string) ([]service_models.MatchStats, error) {
	rawMatches, err := s.repo.GetAllMatchesStats(context, nickname, game)
	if err != nil {
		// Just casual print?
		s.logger.Printf("%v: %v", nickname, err.Error())
		return []service_models.MatchStats{}, err
	}

	var matches = make([]service_models.MatchStats, len(rawMatches))

	for i := range rawMatches {
		matches[i] = convertors.ConvertMatchStats(rawMatches[i])
	}

	sort.Sort(timesort.ByTime(matches))

	return matches, nil
}

func (s service) GetMatchWithMostKills(ctx context.Context, matches []service_models.MatchStats) (service_models.MatchStats, error) {
	var thatMatch service_models.MatchStats
	var kills uint8
	for _, match := range matches {
		if match.Kills > kills {
			kills = match.Kills
			thatMatch = match
		}
	}

	return thatMatch, nil
}

func (s service) GetLongestMatch(ctx context.Context, matches []service_models.MatchStats) (service_models.MatchStats, error) {
	var thatMatch service_models.MatchStats
	var roundsSum uint8
	for _, match := range matches {
		if match.Score[0]+match.Score[1] > roundsSum {
			roundsSum = match.Score[0] + match.Score[1]
			thatMatch = match
		}
	}

	return thatMatch, nil
}

func (s service) GetGlobalKD(ctx context.Context, matches []service_models.MatchStats) (float64, error) {
	var kdSum float64
	for _, match := range matches {
		kdSum += float64(match.Kills) / float64(match.Deaths)
	}

	return kdSum / float64(len(matches)), nil
}
