package val

import (
	"encoding/json"
	"fmt"
	"rito-news/lib/contentstack"
	"rito-news/lib/utils"
	"time"
)

type VALORANTEsportsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type valorantEsportsAPIResponseEntry struct {
	UID     string `json:"uid"`
	Authors []struct {
		Title string `json:"title"`
	} `json:"authors"`
	BannerSettings struct {
		Banner struct {
			Url string `json:"url"`
		} `json:"banner"`
	} `json:"banner_settings"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
	VideoLink string `json:"video_link"`
}

type VALORANTEsports struct {
	Locale string
}

func (VALORANTEsports) getContentStackKeys(params contentstack.ContentStackQueryParameters) *contentstack.ContentStackKeys {
	return contentstack.GetContentStackKeys("https://valorantesports.com/news", "body", &params)
}

func (client VALORANTEsports) getContentStackParameters(count int) contentstack.ContentStackQueryParameters {
	return contentstack.ContentStackQueryParameters{
		ContentType: "articles",
		Locale:      client.Locale,
		Count:       count,
		Environment: "production",
		Filters: map[string][]string{
			"query": {`{"hide_from_newsfeeds":{"$ne":true}}`},
			"only[BASE][]": {
				"title",
				"_content_type_uid",
				"banner_settings",
				"authors",
				"date",
				"description",
				"external_link",
				"url",
				"video_link",
			},
			"include[]": {
				"authors",
			},
			"only[authors][]": {
				"title",
			},
		},
	}
}

func (client VALORANTEsports) getContentStackItems(count int) ([]valorantEsportsAPIResponseEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]valorantEsportsAPIResponseEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (client VALORANTEsports) generateNewsLink(entry valorantEsportsAPIResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.VideoLink != "" {
		return entry.VideoLink
	}
	return fmt.Sprintf("https://valorantesports.com/%s/%s", utils.TrimSlashes(entry.Url.Url), client.Locale)
}

func (client VALORANTEsports) GetItems(count int) ([]VALORANTEsportsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]VALORANTEsportsEntry, len(items))

	for i, item := range items {
		url := client.generateNewsLink(item)

		authors := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authors[i] = author.Title
		}

		results[i] = VALORANTEsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.BannerSettings.Banner.Url,
			Title:       item.Title,
			Url:         url,
		}
	}

	return results, nil
}
