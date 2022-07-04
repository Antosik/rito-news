package lol_source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LeagueOfLegendsNewsEntry struct {
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

type LeagueOfLegendsNewsResponse struct {
	Result struct {
		Data struct {
			AllArticles struct {
				Edges []LeagueOfLegendsNewsEntry `json:"edges"`
			} `json:"allArticles"`
		} `json:"data"`
	} `json:"result"`
}

type LeagueOfLegendsNews struct {
	Locale string
}

func (client LeagueOfLegendsNews) loadItems(count int) ([]LeagueOfLegendsNewsEntry, error) {
	url := fmt.Sprintf(
		"https://www.leagueoflegends.com/page-data/%s/latest-news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response LeagueOfLegendsNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllArticles.Edges[:count], nil
}

func (client LeagueOfLegendsNews) generateNewsLink(entry LeagueOfLegendsNewsEntry) string {
	if entry.Node.YouTubeLink != "" {
		return entry.Node.YouTubeLink
	}
	return fmt.Sprintf("https://www.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Node.Url.Url))
}

func (client LeagueOfLegendsNews) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return nil, fmt.Errorf("can't generate UUID: %w", err)
		}

		authors := make([]string, len(item.Node.Author))
		for i, author := range item.Node.Author {
			authors[i] = author.Title
		}

		var category string
		if len(item.Node.Category) > 0 {
			category = item.Node.Category[0].Title
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Node.Title,
			Summary:   item.Node.Description,
			Url:       url,
			Author:    strings.Join(authors, ","),
			Image:     item.Node.Banner.Url,
			Category:  category,
			CreatedAt: item.Node.Date,
			UpdatedAt: item.Node.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
