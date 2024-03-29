package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
)

// League of Legends news entry
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
	Node struct {
		Author []struct {
			Title string `json:"title"`
		} `json:"author"`
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
		UID          string    `json:"uid"`
		URL          struct {
			URL string `json:"url"`
		} `json:"url"`
		YouTubeLink string `json:"youtube_link"`
	} `json:"node"`
}

type rawNewsResponse struct {
	Result struct {
		Data struct {
			AllArticles struct {
				Edges []rawNewsEntry `json:"edges"`
			} `json:"allArticles"`
		} `json:"data"`
	} `json:"result"`
}

// A client that allows to get official League of Legends news.
//
// Source - https://www.leagueoflegends.com/en-us/news
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, de-de, es-es, fr-fr, it-it, en-pl, pl-pl, el-gr, ro-ro,
	// hu-hu, cs-cz, es-mx, pt-br, ja-jp, ru-ru, tr-tr, en-au, ko-kr,
	// en-sg, en-ph, vi-vn, th-th, zh-tw
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://www.leagueoflegends.com/page-data/%s/latest-news/page-data.json",
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

	edges := response.Result.Data.AllArticles.Edges
	sliceSize := utils.MinInt(count, len(edges))

	return edges[:sliceSize], nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.Node.ExternalLink != "" {
		return entry.Node.ExternalLink
	}

	if entry.Node.YouTubeLink != "" {
		return entry.Node.YouTubeLink
	}

	return fmt.Sprintf("https://www.leagueoflegends.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Node.URL.URL))
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]NewsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		authors := make([]string, len(item.Node.Author))
		for i, author := range item.Node.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Node.Category))
		for i, category := range item.Node.Category {
			categories[i] = category.Title
		}

		tags := make([]string, len(item.Node.ArticleTags))
		for i, tag := range item.Node.ArticleTags {
			tags[i] = tag.Title
		}

		results[i] = NewsEntry{
			UID:         item.Node.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Node.Date,
			Description: item.Node.Description,
			Image:       item.Node.Banner.URL,
			Tags:        tags,
			Title:       item.Node.Title,
			URL:         url,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}
