package aoc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tfaughnan/artoo/config"
)

type response struct {
	Members map[string]member `json:"members"`
	OwnerId int               `json:"owner_id"`
	Event   string            `json:"event"`
}

type member struct {
	Id         int                              `json:"id"`
	Name       string                           `json:"name"`
	LocalScore int                              `json:"local_score"`
	Stars      int                              `json:"stars"`
	Days       map[string](map[string]struct{}) `json:"completion_day_level"`
}

func fetchLeaderboard(cfg config.AocConfig, timeout int) (response, error) {
	if cfg.LeaderboardURL == "" {
		return response{}, errors.New("Missing aoc leaderboard_url in config")
	} else if cfg.SessionCookie == "" {
		return response{}, errors.New("Missing aoc session_cookie in config")
	}

	u := cfg.LeaderboardURL + ".json"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return response{}, err
	}
	req.Header.Set("Cookie", fmt.Sprintf("session=%s", cfg.SessionCookie))
	req.Header.Set("User-Agent", "https://github.com/tfaughnan/artoo")

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("Received status \"%s\"", resp.Status))
		return response{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var r response
	if err := json.Unmarshal(body, &r); err != nil {
		return response{}, err
	}

	return r, nil
}
