package val

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/lib/utils"
	"time"
)

type VALORANTNewsEntry struct {
	UID         string    `json:"uid"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type valorantNewsAPIResponseEntry struct {
	UID    string `json:"uid"`
	Banner struct {
		Url string `json:"url"`
	} `json:"banner"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
}

type valorantNewsAPIResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Nodes []valorantNewsAPIResponseEntry `json:"nodes"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type VALORANTNews struct {
	Locale string
}

func (client VALORANTNews) loadItems(count int) ([]valorantNewsAPIResponseEntry, error) {
	url := fmt.Sprintf(
		"https://playvalorant.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response valorantNewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllContentstackArticles.Nodes[:count], nil
}

func (client VALORANTNews) generateNewsLink(entry valorantNewsAPIResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://playvalorant.com/%s/%s", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client VALORANTNews) GetItems(count int) ([]VALORANTNewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]VALORANTNewsEntry, len(items))

	for i, item := range items {
		url := client.generateNewsLink(item)

		results[i] = VALORANTNewsEntry{
			UID:         item.UID,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.Banner.Url,
			Title:       item.Title,
			Url:         url,
		}
	}

	return results, nil
}
