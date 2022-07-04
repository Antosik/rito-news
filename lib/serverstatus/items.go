package serverstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ServerStatusEntry struct {
	UID         string    `json:"uid"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type serverStatusEntryAPITranslation struct {
	Content string `json:"content"`
	Locale  string `json:"locale"`
}

type serverStatusEntryAPIUpdate struct {
	UpdatedAt    time.Time                         `json:"updated_at"`
	Publish      bool                              `json:"publish"`
	Translations []serverStatusEntryAPITranslation `json:"translations"`
	CreatedAt    time.Time                         `json:"created_at"`
	Author       string                            `json:"author"`
	Id           int                               `json:"id"`
}

type serverStatusAPIEntry struct {
	Titles  []serverStatusEntryAPITranslation `json:"titles"`
	Updates []serverStatusEntryAPIUpdate      `json:"updates"`
}

type serverStatusAPIResponse struct {
	Maintenances []serverStatusAPIEntry `json:"maintenances"`
	Incidents    []serverStatusAPIEntry `json:"incidents"`
}

func GetServerStatusItems(url string, locale string) ([]ServerStatusEntry, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	var response serverStatusAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return TransformServerStatusToNewsItems(response, locale), nil
}
