package wr_source

import (
	"encoding/json"
	"fmt"
	"rito-news/sources/base/contentstack"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WildRiftEsportsEntry struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Authors     []struct {
		Title string `json:"title"`
	} `json:"authors"`
	BannerSettings struct {
		Banner struct {
			Url string `json:"url"`
		} `json:"banner"`
	} `json:"banner_settings"`
	Categories []struct {
		MachineName string `json:"machine_name"`
		Title       string `json:"title"`
	} `json:"category"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
}

type WildRiftEsports struct {
	Locale string
}

func (WildRiftEsports) getContentStackKeys(params contentstack.ContentStackQueryParameters) *contentstack.ContentStackKeys {
	return contentstack.GetContentStackKeys("https://wildriftesports.com/en-us/news", `a[href^="/en-us/news/"]`, &params)
}

func (client WildRiftEsports) getContentStackParameters(count int) contentstack.ContentStackQueryParameters {
	return contentstack.ContentStackQueryParameters{
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
				"article_tags",
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

func (client WildRiftEsports) getContentStackItems(count int) ([]WildRiftEsportsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return []WildRiftEsportsEntry{}, err
	}

	items := make([]WildRiftEsportsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return []WildRiftEsportsEntry{}, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (client WildRiftEsports) generateNewsLink(entry WildRiftEsportsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://wildriftesports.com/%s/%s", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client WildRiftEsports) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.getContentStackItems(count)
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

		authors := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authors[i] = author.Title
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Description,
			Url:       url,
			Author:    strings.Join(authors, ","),
			Image:     item.BannerSettings.Banner.Url,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
