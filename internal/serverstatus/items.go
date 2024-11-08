// Package serverstatus implements the way to interact with Server Status system (https://status.riotgames.com)
package serverstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
)

type Entry struct {
	UID         string    `json:"uid"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawTranslationEntry struct {
	Content string `json:"content"`
	Locale  string `json:"locale"`
}

type rawUpdateEntry struct {
	UpdatedAt    time.Time             `json:"updated_at"`
	Publish      bool                  `json:"publish"`
	Translations []rawTranslationEntry `json:"translations"`
	CreatedAt    time.Time             `json:"created_at"`
	Author       string                `json:"author"`
	ID           int                   `json:"id"`
}

type rawEntry struct {
	Titles  []rawTranslationEntry `json:"titles"`
	Updates []rawUpdateEntry      `json:"updates"`
}

type rawResponse struct {
	Maintenances []rawEntry `json:"maintenances"`
	Incidents    []rawEntry `json:"incidents"`
}

func GetItems(url string, locale string) ([]Entry, error) {
	req, err := utils.NewGETJSONRequest(url)
	if err != nil {
		return nil, fmt.Errorf("can't get json content: %w", err)
	}

	client := &http.Client{} //nolint:exhaustruct

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	var response rawResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return transformRawResponseToEntry(response, locale), nil
}
