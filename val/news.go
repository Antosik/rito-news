package val

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
)

// VALORANT news entry
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

// A client that allows to get official VALORANT news.
//
// Source - https://playvalorant.com/en-us/news/
type NewsClient struct {
	// Available locales:
	// en-US, en-GB, de-DE, es-ES, fr-FR, it-IT, pl-PL, ru-RU, tr-TR,
	// es-MX, id-ID, ja-JP, ko-KR, pt-BR, th-TH, vi-VN, zh-TW, ar-AE
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://playvalorant.com/page-data/%s/news/page-data.json",
		client.Locale,
	)

	req, err := utils.NewGETJSONRequest(url)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response rawNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	nodes := response.Result.Data.AllContentstackArticles.Nodes
	sliceSize := utils.MinInt(count, len(nodes))

	return nodes[:sliceSize], nil
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
