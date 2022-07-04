package tft_source

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TeamfightTacticsNewsEntry struct {
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
}

type TeamfightTacticsNewsResponse struct {
	Result struct {
		Data struct {
			All struct {
				Edges []struct {
					Node struct {
						Entries []TeamfightTacticsNewsEntry `json:"entries"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"all"`
		} `json:"data"`
	} `json:"result"`
}

type TeamfightTacticsNews struct {
	Locale string
}

func (client TeamfightTacticsNews) loadItems(count int) ([]TeamfightTacticsNewsEntry, error) {
	url := fmt.Sprintf(
		"https://teamfighttactics.leagueoflegends.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return []TeamfightTacticsNewsEntry{}, errors.New("Can't load news: " + err.Error())
	}
	defer res.Body.Close()

	var response TeamfightTacticsNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []TeamfightTacticsNewsEntry{}, errors.New("Can't decode response: " + err.Error())
	}

	return response.Result.Data.All.Edges[0].Node.Entries[:count], nil
}

func (client TeamfightTacticsNews) generateNewsLink(entry TeamfightTacticsNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}
	return fmt.Sprintf("https://teamfighttactics.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client TeamfightTacticsNews) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return []abstract.NewsItem{}, errors.New("Can't generate UUID: " + err.Error())
		}

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		var category string
		if len(item.Category) > 0 {
			category = item.Category[0].Title
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Description,
			Url:       url,
			Author:    strings.Join(authors, ","),
			Image:     item.Banner.Url,
			Category:  category,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
