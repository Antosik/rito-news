package wr

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

// Wild Rift news entry
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
	Title  string    `json:"title"`
	Date   time.Time `json:"publishedAt"`
	Action struct {
		Type    string `json:"type"` // 'weblink', 'youtube_video'
		Payload struct {
			URL string `json:"url"`
		} `json:"payload"`
	} `json:"action"`
	Media struct {
		Type string `json:"type"` // 'image'
		URL  string `json:"url"`
	} `json:"media"`
	Description struct {
		Type string `json:"type"` // 'html'
		Body string `json:"body"`
	} `json:"description"`
	Category struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		MachineName string `json:"machineName"`
	} `json:"category"`
}

type rawDataBladeResponse struct {
	Type       string         `json:"type"`
	FragmentID string         `json:"fragmentId"` // should be 'news'
	Items      []rawNewsEntry `json:"items"`
}

type rawDataResponse struct {
	Props struct {
		PageProps struct {
			Page struct {
				Blades []rawDataBladeResponse `json:"blades"`
			} `json:"page"`
		} `json:"pageProps"`
	} `json:"props"`
}

// A client that allows to get official Wild Rift news.
//
// Source - https://wildrift.leagueoflegends.com/en-us/news/
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, fr-fr, de-de, es-es, it-it, pl-pl, ru-ru, tr-tr, id-id,
	// ms-my, pt-br, ja-jp, ko-kr, zh-tw, th-th, vi-vn, es-mx, en-sg, ar-ae
	Locale string
}

func (client NewsClient) loadItems(count int) ([]rawNewsEntry, error) {
	url := fmt.Sprintf(
		"https://wildrift.leagueoflegends.com/%s/news/",
		client.Locale,
	)

	body, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return nil, err
	}

	doc, err := utils.ReadHTML(body)
	if err != nil {
		return nil, fmt.Errorf("can't read news page content: %w", err)
	}

	content := doc.Find("#__NEXT_DATA__").Text()

	var response rawDataResponse
	if err := json.NewDecoder(strings.NewReader(content)).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	for _, item := range response.Props.PageProps.Page.Blades {
		if item.FragmentID != "news" || len(item.Items) == 0 {
			continue
		}

		items := item.Items

		// Sort in case of same publish date
		sort.Slice(items, func(i, j int) bool {
			datecomp := items[i].Date.Compare(items[j].Date)

			if datecomp == 0 {
				return items[i].Title < items[j].Title
			}

			return datecomp > 0
		})

		sliceSize := utils.MinInt(count, len(item.Items))

		return items[:sliceSize], nil
	}

	return nil, fmt.Errorf("can't find news data: %w", err)
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	link := entry.Action.Payload.URL

	switch linkType := entry.Action.Type; linkType {
	case "weblink":
		if strings.HasPrefix(link, "http") {
			return link
		}

		return fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/", utils.TrimSlashes(entry.Action.Payload.URL))
	case "youtube_video":
		return link
	default:
		return fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/news", client.Locale)
	}
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

		authors := make([]string, 0)
		categories := []string{item.Category.Title}
		tags := make([]string, 0)

		results[i] = NewsEntry{
			UID:         uid,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description.Body,
			Image:       item.Media.URL,
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
