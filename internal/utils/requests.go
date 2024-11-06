package utils

import (
	"fmt"
	"io"
	"net/http"
)

const UserAgent = "Antosik/rito-news (https://github.com/Antosik/rito-news)"

func NewGETRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	return req, nil
}

func NewGETJSONRequest(url string) (*http.Request, error) {
	req, err := NewGETRequest(url)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}

func RunGETHTMLRequest(url string) (string, error) {
	req, err := NewGETRequest(url)
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{} //nolint:exhaustruct

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unsuccessful request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("can't decode response: %w", err)
	}

	return string(body), nil
}
