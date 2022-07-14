package tft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
)

// Teamfight Tactics news entry
type NewsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawNewsEntry struct {
	UID    string `json:"uid"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Banner struct {
		URL string `json:"url"`
	} `json:"banner"`
	Category []struct {
		Title string `json:"title"`
	} `json:"category"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	URL          struct {
		URL string `json:"url"`
	} `json:"url"`
	YouTubeLink string `json:"youtube_link"`
}

type teamfightTacticsNewsAPIResponse struct {
	Result struct {
		Data struct {
			All struct {
				Edges []struct {
					Node struct {
						Entries []rawNewsEntry `json:"entries"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"all"`
		} `json:"data"`
	} `json:"result"`
}

// A client that allows to get official Teamfight Tactics news.
//
// Source - https://teamfighttactics.leagueoflegends.com/
type NewsClient struct {
	// Available locales:
	// en-US, en-GB, de-DE, es-ES, fr-FR, it-IT, en-AU, pl-PL, ru-RU,
	// el-GR, ro-RO, hu-HU, cs-CZ, es-MX, pt-BR, tr-TR, ko-KR, ja-JP
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://teamfighttactics.leagueoflegends.com/page-data/%s/news/page-data.json",
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

	var response teamfightTacticsNewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	entries := response.Result.Data.All.Edges[0].Node.Entries
	sliceSize := utils.MinInt(count, len(entries))

	return entries[:sliceSize], nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}

	return fmt.Sprintf(
		"https://teamfighttactics.leagueoflegends.com/%s/%s/",
		client.Locale,
		utils.TrimSlashes(entry.URL.URL),
	)
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]NewsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Category))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		results[i] = NewsEntry{
			UID:         item.UID,
			Authors:     authors,
			Categories:  categories,
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
