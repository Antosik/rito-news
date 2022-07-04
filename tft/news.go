package tft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/lib/utils"
	"time"
)

type TeamfightTacticsNewsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type teamfightTacticsNewsAPIResponseEntry struct {
	UID    string `json:"uid"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Banner struct {
		Url string `json:"url"`
	} `json:"banner"`
	Category []struct {
		Title string `json:"title"`
	} `json:"category"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
	YouTubeLink string `json:"youtube_link"`
}

type teamfightTacticsNewsAPIResponse struct {
	Result struct {
		Data struct {
			All struct {
				Edges []struct {
					Node struct {
						Entries []teamfightTacticsNewsAPIResponseEntry `json:"entries"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"all"`
		} `json:"data"`
	} `json:"result"`
}

type TeamfightTacticsNews struct {
	Locale string
}

func (client TeamfightTacticsNews) loadItems(count int) ([]teamfightTacticsNewsAPIResponseEntry, error) {
	url := fmt.Sprintf(
		"https://teamfighttactics.leagueoflegends.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response teamfightTacticsNewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.All.Edges[0].Node.Entries[:count], nil
}

func (client TeamfightTacticsNews) generateNewsLink(entry teamfightTacticsNewsAPIResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}
	return fmt.Sprintf("https://teamfighttactics.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client TeamfightTacticsNews) GetItems(count int) ([]TeamfightTacticsNewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]TeamfightTacticsNewsEntry, len(items))

	for i, item := range items {
		url := client.generateNewsLink(item)

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Author))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		results[i] = TeamfightTacticsNewsEntry{
			UID:         item.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.Banner.Url,
			Title:       item.Title,
			Url:         url,
		}
	}

	return results, nil
}
