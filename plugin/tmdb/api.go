package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tfaughnan/artoo/config"
)

type searchResponse struct {
	Page         int            `json:"page"`
	Results      []searchResult `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
}

type searchResult struct {
	Adult        bool    `json:"adult"`
	BackdropPath string  `json:"backdrop_path"`
	GenreIDs     []int   `json:"genre_ids"`
	ID           int     `json:"id"`
	OrigLang     string  `json:"original_language"`
	OrigTitle    string  `json:"original_title"`
	Overview     string  `json:"overview"`
	Popularity   float32 `json:"popularity"`
	PosterPath   string  `json:"poster_path"`
	ReleaseDate  string  `json:"release_date"`
	Title        string  `json:"title"`
	Video        bool    `json:"video"`
	VoteAverage  float32 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
}

type creditsResponse struct {
	ID int `json:"id"`
	// XXX: Cast []creditsPerson `json:"cast"`
	Crew []creditsPerson `json:"crew"`
}

type creditsPerson struct {
	Adult        bool    `json:"adult"`
	Gender       int     `json:"gender"`
	ID           int     `json:"id"`
	KnownForDept string  `json:"known_for_department"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Popularity   float32 `json:"popularity"`
	ProfilePath  string  `json:"profile_path"`
	CreditID     string  `json:"credit_id"`
	Dept         string  `json:"department"`
	Job          string  `json:"job"`
}

func fetchSearch(cfg config.TmdbConfig, timeout int, query string) (searchResponse, error) {
	v := url.Values{}
	v.Set("api_key", cfg.ApiKey)
	v.Set("query", query)
	u := fmt.Sprintf("%s/search/movie?%s", cfg.ApiURL, v.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return searchResponse{}, err
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return searchResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("Received status \"%s\"", resp.Status))
		return searchResponse{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return searchResponse{}, err
	}

	var sresp searchResponse
	if err := json.Unmarshal(body, &sresp); err != nil {
		return searchResponse{}, err
	}

	return sresp, nil
}

func fetchCredits(cfg config.TmdbConfig, timeout int, id int) (creditsResponse, error) {
	v := url.Values{}
	v.Set("api_key", cfg.ApiKey)
	u := fmt.Sprintf("%s/movie/%d/credits?%s", cfg.ApiURL, id, v.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return creditsResponse{}, err
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return creditsResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("Received status \"%s\"", resp.Status))
		return creditsResponse{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return creditsResponse{}, err
	}

	var cresp creditsResponse
	if err := json.Unmarshal(body, &cresp); err != nil {
		return creditsResponse{}, err
	}

	return cresp, nil
}