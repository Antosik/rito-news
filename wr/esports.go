package wr

import (
	"encoding/json"
	"fmt"
	"rito-news/lib/contentstack"
	"rito-news/lib/utils"
	"time"
)

type WildRiftEsportsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type wildRiftEsportsApiResponseEntry struct {
	UID     string `json:"uid"`
	Authors []struct {
		Title string `json:"title"`
	} `json:"authors"`
	BannerSettings struct {
		Banner struct {
			Url string `json:"url"`
		} `json:"banner"`
	} `json:"banner_settings"`
	Category []struct {
		Title string `json:"title"`
	} `json:"category"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
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

func (client WildRiftEsports) getContentStackItems(count int) ([]wildRiftEsportsApiResponseEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]wildRiftEsportsApiResponseEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (client WildRiftEsports) generateNewsLink(entry wildRiftEsportsApiResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://wildriftesports.com/%s/%s", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client WildRiftEsports) GetItems(count int) ([]WildRiftEsportsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]WildRiftEsportsEntry, len(items))

	for i, item := range items {
		url := client.generateNewsLink(item)

		authors := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Category))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		results[i] = WildRiftEsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.BannerSettings.Banner.Url,
			Title:       item.Title,
			Url:         url,
		}
	}

	return results, nil
}
