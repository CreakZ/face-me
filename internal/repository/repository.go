package repository

import (
	"context"
	"encoding/json"
	"errors"
	"faceit_parser/internal/constants"
	repository_models "faceit_parser/internal/models/repository"
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
	errs   custom_errors.Errs
}

func InitRepository() repository {
	return repository{
		logger: log.Default(),
		errs:   custom_errors.InitErrs(),
	}
}

type Repository interface {
	GetPlayerID(nickname string) (string, error)
	GetPlayerMatchesCount(nickname, game string) (int, error)
	GetAllMatchesStats(context context.Context, nickname, game string) ([]repository_models.MatchStatsRaw, error)
}

// GetPlayerID returns player ID according to his nickname
// REMOVE THIS FUNCTION SOMEWHERE
func (r repository) GetPlayerID(nickname string) (string, error) {
	URL := strings.Join([]string{constants.BaseURL, constants.IDInfo, nickname}, "/")

	// request to get player ID
	req, err := http.Get(URL)
	if err != nil {
		return "", err
	}

	if req.StatusCode == http.StatusNotFound {
		return "", errors.New("wrong nickname")
	}

	body, _ := io.ReadAll(req.Body)

	type Payload struct {
		Data map[string]interface{} `json:"payload"`
	}

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	return payload.Data["id"].(string), nil
}

// GetPlayerMatches returns number of matches user with specified id played in needed game
// REMOVE THIS FUNCTION SOMEWHERE
func (r repository) GetPlayerMatchesCount(nickname, game string) (int, error) {
	id, err := r.GetPlayerID(nickname)
	if err != nil {
		return 0, err
	}

	URL := strings.Join([]string{constants.BaseURL, fmt.Sprintf(constants.MatchInfoURL, id, game)}, "/")

	req, err := http.Get(URL)
	if err != nil {
		return 0, err
	}

	body, _ := io.ReadAll(req.Body)

	type Stats struct {
		Lifetime map[string]interface{} `json:"lifetime"`
	}

	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		return 0, err
	}

	return int(stats.Lifetime["rev"].(float64)), nil
}

func (r repository) GetAllMatchesStats(context context.Context, nickname, game string) ([]repository_models.MatchStatsRaw, error) {
	id, err := r.GetPlayerID(nickname)
	if err != nil {
		return []repository_models.MatchStatsRaw{}, err
	}

	matchCount, err := r.GetPlayerMatchesCount(nickname, game)
	if err != nil {
		return []repository_models.MatchStatsRaw{}, err
	}

	fmt.Println(matchCount)

	var page, curr int
	bodies := make(chan []byte, matchCount/constants.BaseSize+1)
	errs, _ := errgroup.WithContext(context)

	for matchCount > curr {
		URL := strings.Join([]string{constants.BaseURL,
			fmt.Sprintf(constants.MatchInfo,
				id,
				game,
				page,
				constants.BaseSize)},
			"/")

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
		r.logger.Printf("no matches played in %v\n", game)
		return []repository_models.MatchStatsRaw{}, r.errs.NoMatches(game)
	}

	if count != matchCount {
		r.logger.Printf("info about first %v matches failed to fetch due to external reasons\n", matchCount-len(allMatches))
	}

	return allMatches, nil
}
