package contentstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ContentStackQueryParameters struct {
	ContentType string
	Locale      string
	Count       int
	Environment string
	Filters     map[string][]string
}

func GetContentStackItems(keys *ContentStackKeys, parameters *ContentStackQueryParameters) ([]json.RawMessage, error) {
	if keys.access_token == "" || keys.api_key == "" {
		return []json.RawMessage{}, fmt.Errorf("incorrect api keys: " + keys.String())
	}

	url := generateContentStackUrl(parameters)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("api_key", keys.api_key)
	req.Header.Set("access_token", keys.access_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []json.RawMessage{}, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	response := struct {
		Entries []json.RawMessage `json:"entries"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return []json.RawMessage{}, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Entries, nil
}

func generateContentStackUrl(parameters *ContentStackQueryParameters) string {
	var query []string

	for key, values := range parameters.Filters {
		for _, val := range values {
			query = append(query, url.PathEscape(fmt.Sprintf("%s=%s", key, val)))
		}
	}

	return fmt.Sprintf(`https://cdn.contentstack.io/v3/content_types/%s/entries/?locale=%s&environment=%s&limit=%d&desc=date&%v`,
		parameters.ContentType,
		parameters.Locale,
		parameters.Environment,
		parameters.Count,
		strings.Join(query, "&"),
	)
}
