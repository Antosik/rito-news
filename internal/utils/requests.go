package utils

import (
	"fmt"
	"net/http"
)

const UserAgent = "Antosik/rito-news (https://github.com/Antosik/rito-news)"

func NewGETJSONRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("user-agent", UserAgent)

	return req, err
}
