package service

import (
	"context"
	"faceit_parser/internal/convertors"
	service_models "faceit_parser/internal/models/service"
	"faceit_parser/internal/repository"
	"faceit_parser/pkg/timesort"
	"log"
	"sort"

	"github.com/redis/go-redis/v9"
)

type service struct {
	repo   repository.Repository
	logger *log.Logger
	cache  *redis.Client
}

func InitService(repo repository.Repository, logger *log.Logger, cache *redis.Client) service {
	return service{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

type Service interface {
	// pagination need to be done
	// GetMatchesStats(context context.Context, nickname, game string, page, size int) ([]service_models.MatchStats, error)
	GetAllMatchesStats(context context.Context, nickname, game string) ([]service_models.MatchStats, error)
	GetMatchWithMostKills(ctx context.Context, sessionID string) (service_models.MatchStats, error)
	GetLongestMatch(ctx context.Context, sessionID string) (service_models.MatchStats, error)
	GetGlobalKD(ctx context.Context, sessionID string) (float64, error)
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

// TODO
func (s service) GetMatchWithMostKills(ctx context.Context, sessionID string) (service_models.MatchStats, error) {
	return service_models.MatchStats{}, nil
}

// TODO
func (s service) GetLongestMatch(ctx context.Context, sessionID string) (service_models.MatchStats, error) {
	return service_models.MatchStats{}, nil
}

func (s service) GetGlobalKD(ctx context.Context, sessionID string) (float64, error) {
	return 0.0, nil
}
