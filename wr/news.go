package wr

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
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawNewsEntry struct {
	UID        string `json:"id"`
	Categories []struct {
		Title string `json:"title"`
	} `json:"categories"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	ExternalLink  string    `json:"externalLink"`
	FeaturedImage struct {
		Banner struct {
			URL string `json:"url"`
		} `json:"banner"`
	} `json:"featuredImage"`
	Link struct {
		URL string `json:"url"`
	} `json:"link"`
	Tags []struct {
		Title string `json:"title"`
	} `json:"tags"`
	Title       string `json:"title"`
	YouTubeLink string `json:"youtubeLink"`
}

type rawNewsResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Articles []rawNewsEntry `json:"articles"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type NewsClient struct {
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://wildrift.leagueoflegends.com/page-data/%s/news/page-data.json",
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

	return response.Result.Data.AllContentstackArticles.Articles[:count], nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}

	return fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Link.URL))
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	items := make([]NewsEntry, len(stackItems))

	for i, item := range stackItems {
		url := client.getLinkForEntry(item)

		categories := make([]string, len(item.Categories))
		for i, category := range item.Categories {
			categories[i] = category.Title
		}

		tags := make([]string, len(item.Tags))
		for i, tag := range item.Tags {
			tags[i] = tag.Title
		}

		items[i] = NewsEntry{
			UID:         item.UID,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.FeaturedImage.Banner.URL,
			Tags:        tags,
			Title:       item.Title,
			URL:         url,
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	return items, nil
}
