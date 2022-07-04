package wr_source

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

type WildRiftNewsEntry struct {
	Categories []struct {
		Title string `json:"title"`
	} `json:"categories"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	ExternalLink  string    `json:"externalLink"`
	FeaturedImage struct {
		Banner struct {
			Url string `json:"url"`
		} `json:"banner"`
	} `json:"featuredImage"`
	Link struct {
		Url string `json:"url"`
	} `json:"link"`
	Tags []struct {
		Title string `json:"title"`
	} `json:"tags"`
	Title       string `json:"title"`
	UID         string `json:"uid"`
	YouTubeLink string `json:"youtubeLink"`
}

type WildRiftNewsResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Articles []WildRiftNewsEntry `json:"articles"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type WildRiftNews struct {
	Locale string
}

func (client WildRiftNews) loadItems(count int) ([]WildRiftNewsEntry, error) {
	url := fmt.Sprintf(
		"https://wildrift.leagueoflegends.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return []WildRiftNewsEntry{}, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response WildRiftNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []WildRiftNewsEntry{}, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllContentstackArticles.Articles[:count], nil
}

func (client WildRiftNews) generateNewsLink(entry WildRiftNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}
	return fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Link.Url))
}

func (client WildRiftNews) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return []abstract.NewsItem{}, fmt.Errorf("can't generate UUID: %w", err)
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Description,
			Url:       url,
			Image:     item.FeaturedImage.Banner.Url,
			Category:  item.Categories[0].Title,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
