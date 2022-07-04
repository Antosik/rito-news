package serverstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ServerStatusEntryTranslation struct {
	Content string `json:"content"`
	Locale  string `json:"locale"`
}

type ServerStatusEntryUpdate struct {
	UpdatedAt    time.Time                      `json:"updated_at"`
	Publish      bool                           `json:"publish"`
	Translations []ServerStatusEntryTranslation `json:"translations"`
	CreatedAt    time.Time                      `json:"created_at"`
	Author       string                         `json:"author"`
	Id           int                            `json:"id"`
}

type ServerStatusEntry struct {
	Titles  []ServerStatusEntryTranslation `json:"titles"`
	Updates []ServerStatusEntryUpdate      `json:"updates"`
}

type ServerStatusResponse struct {
	Maintenances []ServerStatusEntry `json:"maintenances"`
	Incidents    []ServerStatusEntry `json:"incidents"`
}

func GetServerStatusItems(url string) (ServerStatusResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ServerStatusResponse{}, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ServerStatusResponse{}, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	var response ServerStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ServerStatusResponse{}, fmt.Errorf("can't decode response: %w", err)
	}

	return response, nil
}
