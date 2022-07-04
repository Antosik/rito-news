package wr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/lib/utils"
	"time"
)

type WildRiftNewsEntry struct {
	UID         string    `json:"uid"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type wildRiftNewsAPIResponseEntry struct {
	UID        string `json:"id"`
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
	YouTubeLink string `json:"youtubeLink"`
}

type wildRiftNewsAPIResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Articles []wildRiftNewsAPIResponseEntry `json:"articles"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type WildRiftNews struct {
	Locale string
}

func (client WildRiftNews) loadItems(count int) ([]wildRiftNewsAPIResponseEntry, error) {
	url := fmt.Sprintf(
		"https://wildrift.leagueoflegends.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response wildRiftNewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllContentstackArticles.Articles[:count], nil
}

func (client WildRiftNews) generateNewsLink(entry wildRiftNewsAPIResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.YouTubeLink != "" {
		return entry.YouTubeLink
	}
	return fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Link.Url))
}

func (client WildRiftNews) GetItems(count int) ([]WildRiftNewsEntry, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	items := make([]WildRiftNewsEntry, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		categories := make([]string, len(item.Categories))
		for i, category := range item.Categories {
			categories[i] = category.Title
		}

		tags := make([]string, len(item.Tags))
		for i, tag := range item.Tags {
			tags[i] = tag.Title
		}

		items[i] = WildRiftNewsEntry{
			UID:         item.UID,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.FeaturedImage.Banner.Url,
			Tags:        tags,
			Title:       item.Title,
			Url:         url,
		}
	}

	return items, nil
}
