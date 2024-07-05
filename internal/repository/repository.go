package repository

import (
	"context"
	"encoding/json"
	"faceit_parser/internal/constants"
	repository_models "faceit_parser/internal/models/repository"
	"faceit_parser/internal/requests"
	custom_errors "faceit_parser/pkg/errs"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"
)

type repository struct {
	logger *log.Logger
}

func InitRepository(logger *log.Logger) repository {
	return repository{
		logger: logger,
	}
}

type Repository interface {
	GetAllMatchesStats(context context.Context, nickname, game string) ([]repository_models.MatchStatsRaw, error)
}

func (r repository) GetAllMatchesStats(context context.Context, nickname, game string) ([]repository_models.MatchStatsRaw, error) {
	id, matchCount, err := requests.GetPlayerData(nickname, game)
	if err != nil {
		return []repository_models.MatchStatsRaw{}, err
	}

	var page, curr int
	bodies := make(chan []byte, matchCount/constants.BaseSize+1)
	errs, _ := errgroup.WithContext(context)

	for matchCount > curr {
		URL := strings.Join([]string{constants.BaseURL, fmt.Sprintf(constants.MatchInfo, id, game, page, constants.BaseSize)}, "/")

		errs.Go(func() error {
			req, err := http.Get(URL)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			body, _ := io.ReadAll(req.Body)
			bodies <- body

			return nil
		})

		page++
		curr += constants.BaseSize
	}

	reqErr := errs.Wait()
	if reqErr != nil {
		return []repository_models.MatchStatsRaw{}, reqErr
	}

	var allMatches, tempMatchesRaw []repository_models.MatchStatsRaw
	for len(bodies) != 0 {
		body := <-bodies

		if unmarshalErr := json.Unmarshal(body, &tempMatchesRaw); unmarshalErr != nil {
			return []repository_models.MatchStatsRaw{}, unmarshalErr
		}

		allMatches = append(allMatches, tempMatchesRaw...)
	}

	var count = len(allMatches)

	if count == 0 {
		return []repository_models.MatchStatsRaw{}, custom_errors.NoMatches()
	}

	if count != matchCount {
		r.logger.Printf("info about first %v matches failed to fetch due to external reasons\n", matchCount-count)
	}

	return allMatches, nil
}
