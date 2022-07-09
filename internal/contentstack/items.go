package contentstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Antosik/rito-news/internal/utils"
)

type Parameters struct {
	ContentType string
	Locale      string
	Count       int
	Environment string
	Filters     map[string][]string
}

func GetItems(keys *Keys, parameters *Parameters) ([]json.RawMessage, error) {
	if keys.accessToken == "" || keys.apiKey == "" {
		return nil, fmt.Errorf("incorrect api keys: %s", keys.String())
	}

	url := generateURL(parameters)

	req, err := utils.NewGETJSONRequest(url)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api_key", keys.apiKey)
	req.Header.Set("access_token", keys.accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	response := struct {
		Entries []json.RawMessage `json:"entries"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Entries, nil
}

func generateURL(parameters *Parameters) string {
	var query []string

	for key, values := range parameters.Filters {
		for _, val := range values {
			query = append(query, url.PathEscape(fmt.Sprintf("%s=%s", key, val)))
		}
	}

	return fmt.Sprintf(
		`https://cdn.contentstack.io/v3/content_types/%s/entries/?locale=%s&environment=%s&limit=%d&desc=date&%v`,
		parameters.ContentType,
		parameters.Locale,
		parameters.Environment,
		parameters.Count,
		strings.Join(query, "&"),
	)
}
