package requests

import (
	"encoding/json"
	"faceit_parser/internal/constants"
	custom_errors "faceit_parser/pkg/errs"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetPlayerData(nickname, game string) (string, int, error) {
	URL := strings.Join([]string{constants.BaseURL, constants.IDInfo, nickname}, "/")

	// request to get player ID
	req, err := http.Get(URL)
	if err != nil {
		return "", 0, err
	}

	if req.StatusCode == http.StatusNotFound {
		return "", 0, custom_errors.WrongNickname()
	}

	body, _ := io.ReadAll(req.Body)

	type Payload struct {
		Data map[string]interface{} `json:"payload"`
	}

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", 0, err
	}

	id := payload.Data["id"].(string)

	URL = strings.Join([]string{constants.BaseURL, fmt.Sprintf(constants.MatchInfoURL, id, game)}, "/")

	req, err = http.Get(URL)
	if err != nil {
		return "", 0, err
	}

	body, _ = io.ReadAll(req.Body)

	type Stats struct {
		Lifetime map[string]interface{} `json:"lifetime"`
	}

	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		return "", 0, err
	}

	if stats.Lifetime["rev"] == nil {
		return "", 0, err
	}

	return id, int(stats.Lifetime["rev"].(float64)), nil
}
