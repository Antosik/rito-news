package lor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

// Legends of Runeterra news entry
type NewsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawNewsEntry struct {
	UID         string `json:"uid"`
	ArticleTags []struct {
		Title string `json:"title"`
	} `json:"article_tags"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Category []struct {
		Title string `json:"title"`
	} `json:"category"`
	CoverImage struct {
		URL string `json:"url"`
	} `json:"cover_image"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	URL          struct {
		URL string `json:"url"`
	} `json:"url"`
}

type rawDataResponse struct {
	List []rawNewsEntry `json:"list"`
}

// A client that allows to get official Legends of Runeterra news
//
// Source - https://playruneterra.com/en-us/news/
type NewsClient struct {
	// Available locales:
	// en-us, ko-kr, fr-fr, es-es, es-mx, de-de, it-it, pl-pl,
	// pt-br, tr-tr, ru-ru, ja-jp, en-sg, zh-tw, th-th, vi-vn
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://playruneterra.com/api/articles?locale=%s&offset=0&limit=%d",
		client.Locale,
		count,
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

	var response rawDataResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	sliceSize := utils.MinInt(count, len(response.List))

	return response.List[:sliceSize], nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	return fmt.Sprintf("https://playruneterra.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.URL.URL))
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]NewsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)
		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Category))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		tags := make([]string, len(item.ArticleTags))
		for i, tag := range item.ArticleTags {
			tags[i] = tag.Title
		}

		results[i] = NewsEntry{
			UID:         uid,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Summary,
			Image:       item.CoverImage.URL,
			Tags:        tags,
			Title:       item.Title,
			URL:         url,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}
