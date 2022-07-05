package wr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Antosik/rito-news/internal/contentstack"
	"github.com/Antosik/rito-news/internal/utils"
)

type EsportsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawEsportsEntry struct {
	UID     string `json:"uid"`
	Authors []struct {
		Title string `json:"title"`
	} `json:"authors"`
	BannerSettings struct {
		Banner struct {
			URL string `json:"url"`
		} `json:"banner"`
	} `json:"banner_settings"`
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

type EsportsClient struct {
	Locale string
}

func (EsportsClient) getContentStackKeys(params contentstack.Parameters) *contentstack.Keys {
	return contentstack.GetKeys("https://wildriftesports.com/en-us/news", `a[href^="/en-us/news/"]`, &params)
}

func (client EsportsClient) getContentStackParameters(count int) contentstack.Parameters {
	return contentstack.Parameters{
		ContentType: "articles",
		Locale:      client.Locale,
		Count:       count,
		Environment: "production",
		Filters: map[string][]string{
			"query": {`{"$and":[{"hide_from_newsfeeds":{"$ne":"true"}},{"article_type":{"$ne":"Team page"}}]}`},
			"only[BASE][]": {
				"title",
				"_content_type_uid",
				"banner_settings",
				"authors",
				"category",
				"date",
				"description",
				"external_link",
				"url",
			},
			"include[]": {
				"authors",
				"category",
			},
			"only[authors][]": {
				"title",
			},
			"only[category][]": {
				"machine_name",
				"title",
				"url",
			},
		},
	}
}

func (client EsportsClient) getContentStackItems(count int) ([]rawEsportsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]rawEsportsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (client EsportsClient) getLinkForEntry(entry rawEsportsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	return fmt.Sprintf("https://wildriftesports.com/%s/%s", client.Locale, utils.TrimSlashes(entry.URL.URL))
}

func (client EsportsClient) GetItems(count int) ([]EsportsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]EsportsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		authors := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Category))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		results[i] = EsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.BannerSettings.Banner.URL,
			Title:       item.Title,
			URL:         url,
		}
	}

	return results, nil
}
