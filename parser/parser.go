package parser

import (
	"context"
	"encoding/json"
	"errors"
	"faceit_parser/constants"
	"faceit_parser/convertors"
	"faceit_parser/models"
	"faceit_parser/timesort"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
)

// GetPlayerID returns player ID according to his nickname
func GetPlayerID(nickname string) (string, error) {
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
func GetPlayerMatchesCount(nickname, game string) (int, error) {
	id, err := GetPlayerID(nickname)
	if err != nil {
		return 0, err
	}

	URL := strings.Join([]string{constants.BaseURL,
		fmt.Sprintf(constants.MatchInfoURL, id, game)},
		"/")

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

func GetAllMatchesStats(context context.Context, nickname, game string) ([]models.MatchStats, error) {
	id, err := GetPlayerID(nickname)
	if err != nil {
		return []models.MatchStats{}, err
	}

	matchCount, err := GetPlayerMatchesCount(nickname, game)
	if err != nil {
		return []models.MatchStats{}, err
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
		return []models.MatchStats{}, reqErr
	}

	var allMatchesRaw, tempMatchesRaw []models.MatchStatsRaw

	for len(bodies) != 0 {
		body := <-bodies

		if unmarshalErr := json.Unmarshal(body, &tempMatchesRaw); unmarshalErr != nil {
			return []models.MatchStats{}, unmarshalErr
		}

		allMatchesRaw = append(allMatchesRaw, tempMatchesRaw...)
	}

	if matchCount != len(allMatchesRaw) {
		log.Printf("info about first %v matches lost due to external reasons\n", matchCount-len(allMatchesRaw))
	}

	var allMatches = make([]models.MatchStats, len(allMatchesRaw))
	for i := range allMatchesRaw {
		allMatches[i] = convertors.ConvertMatchStats(allMatchesRaw[i])
	}

	sort.Sort(timesort.ByTime(allMatches))

	return allMatches, nil
}

func GetMatchWithMostKills(ctx context.Context, nickname string, game string) (models.MatchStats, error) {
	matches, err := GetAllMatchesStats(ctx, nickname, game)
	if err != nil {
		return models.MatchStats{}, err
	}

	var thatMatch models.MatchStats
	var kills uint8
	for _, match := range matches {
		if match.Kills > kills {
			kills = match.Kills
			thatMatch = match
		}
	}

	return thatMatch, nil
}

func GetLongestMatch(ctx context.Context, nickname, game string) (models.MatchStats, error) {

	matches, err := GetAllMatchesStats(ctx, nickname, game)
	if err != nil {
		return models.MatchStats{}, err
	}

	var thatMatch models.MatchStats
	var roundsSum uint8
	for _, match := range matches {
		if match.Score[0]+match.Score[1] > roundsSum {
			roundsSum = match.Score[0] + match.Score[1]
			thatMatch = match
		}
	}

	return thatMatch, nil
}

func GetGlobalKD(ctx context.Context, nickname, game string) (float64, error) {

	matches, err := GetAllMatchesStats(ctx, nickname, game)
	if err != nil {
		return 0.0, err
	}

	if len(matches) == 0 {
		return 0.0, fmt.Errorf("no matches played in %v", game)
	}

	allMatches, err := GetPlayerMatchesCount(nickname, game)
	if err != nil {
		return 0.0, err
	}

	var kdSum float64
	for _, match := range matches {
		kdSum += float64(match.Kills) / float64(match.Deaths)
	}

	return kdSum / float64(allMatches), nil
}
