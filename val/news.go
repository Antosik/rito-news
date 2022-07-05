package val

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
)

type NewsEntry struct {
	UID         string    `json:"uid"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawNewsEntry struct {
	UID    string `json:"uid"`
	Banner struct {
		URL string `json:"url"`
	} `json:"banner"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	URL          struct {
		URL string `json:"url"`
	} `json:"url"`
}

type rawNewsResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Nodes []rawNewsEntry `json:"nodes"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type NewsClient struct {
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://playvalorant.com/page-data/%s/news/page-data.json",
		client.Locale,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response rawNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllContentstackArticles.Nodes[:count], nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	return fmt.Sprintf("https://playvalorant.com/%s/%s", client.Locale, utils.TrimSlashes(entry.URL.URL))
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]NewsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		results[i] = NewsEntry{
			UID:         item.UID,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.Banner.URL,
			Title:       item.Title,
			URL:         url,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}
