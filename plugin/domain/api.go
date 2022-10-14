package domain

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

type domainResponse struct {
	Products []product `json:"products"`
}

type domain400 struct {
	Message string `json:"message"`
}

type product struct {
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Prices []price `json:"prices"`
}

type price struct {
	PriceBeforeTaxes float32 `json:"price_before_taxes"`
}

func fetchDomain(cfg config.DomainConfig, timeout int, query string) (domainResponse, error) {
	v := url.Values{}
	v.Set("name", query)
	v.Set("currency", cfg.Currency)
	u := fmt.Sprintf("%s/domain/check?%s", cfg.ApiURL, v.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return domainResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Apikey %s", cfg.ApiKey))

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return domainResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 400 {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			var d400 domain400
			if err := json.Unmarshal(body, &d400); err != nil {
				return domainResponse{}, err
			}

			err = errors.New(fmt.Sprintf("400 Bad Response: %s", d400.Message))
			return domainResponse{}, err
		}

		err := errors.New(fmt.Sprintf("Received status \"%s\"", resp.Status))
		return domainResponse{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var dr domainResponse
	if err := json.Unmarshal(body, &dr); err != nil {
		return domainResponse{}, err
	}

	return dr, nil
}
