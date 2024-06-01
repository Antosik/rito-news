package val

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

// VALORANT news entry
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
	UID         string `json:"uid"`
	ArticleTags []struct {
		Title string `json:"title"`
	} `json:"article_tags"`
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
	// en-us, en-gb, de-de, es-es, fr-fr, it-it, pl-pl, ru-ru, tr-tr,
	// es-mx, id-id, ja-jp, ko-kr, pt-br, th-th, vi-vn, zh-tw, ar-ae
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
		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

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
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.Banner.URL,
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
