package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/lib/utils"
	"time"
)

type LeagueOfLegendsNewsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type leagueOfLegendsNewsAPIResponseEntry struct {
	Node struct {
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
		UID          string    `json:"uid"`
		Url          struct {
			Url string `json:"url"`
		} `json:"url"`
		YouTubeLink string `json:"youtube_link"`
	} `json:"node"`
}

type leagueOfLegendsNewsAPIResponse struct {
	Result struct {
		Data struct {
			AllArticles struct {
				Edges []leagueOfLegendsNewsAPIResponseEntry `json:"edges"`
			} `json:"allArticles"`
		} `json:"data"`
	} `json:"result"`
}

type LeagueOfLegendsNews struct {
	Locale string
}

func (client LeagueOfLegendsNews) loadItems(count int) ([]leagueOfLegendsNewsAPIResponseEntry, error) {
	url := fmt.Sprintf(
		"https://www.leagueoflegends.com/page-data/%s/latest-news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response leagueOfLegendsNewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllArticles.Edges[:count], nil
}

func (client LeagueOfLegendsNews) generateNewsLink(entry leagueOfLegendsNewsAPIResponseEntry) string {
	if entry.Node.YouTubeLink != "" {
		return entry.Node.YouTubeLink
	}
	return fmt.Sprintf("https://www.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Node.Url.Url))
}

func (client LeagueOfLegendsNews) GetItems(count int) ([]LeagueOfLegendsNewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]LeagueOfLegendsNewsEntry, len(items))
	for i, item := range items {
		url := client.generateNewsLink(item)

		authors := make([]string, len(item.Node.Author))
		for i, author := range item.Node.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Node.Category))
		for i, category := range item.Node.Category {
			categories[i] = category.Title
		}

		results[i] = LeagueOfLegendsNewsEntry{
			UID:         item.Node.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Node.Date,
			Description: item.Node.Description,
			Image:       item.Node.Banner.Url,
			Title:       item.Node.Title,
			Url:         url,
		}
	}

	return results, nil
}
