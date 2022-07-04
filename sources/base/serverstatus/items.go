package serverstatus

import (
	"encoding/json"
	"errors"
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
	Id           uint16                         `json:"id"`
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
		return ServerStatusResponse{}, errors.New("Can't create request: " + err.Error())
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ServerStatusResponse{}, errors.New("Unsuccessful request: " + err.Error())
	}
	defer resp.Body.Close()

	var response ServerStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ServerStatusResponse{}, errors.New("Can't decode response: " + err.Error())
	}

	return response, nil
}
