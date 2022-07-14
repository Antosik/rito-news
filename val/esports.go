package val

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/contentstack"
	"github.com/Antosik/rito-news/internal/utils"
)

// VALORANT esports news entry
type EsportsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
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
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	URL          struct {
		URL string `json:"url"`
	} `json:"url"`
	VideoLink string `json:"video_link"`
}

// A client that allows to get VALORANT esports news.
//
// Source - https://valorantesports.com/news
type EsportsClient struct {
	// Available locales:
	// en-US, en-GB, en-AU, de-DE, es-ES, es-MX, fr-FR, it-IT, pl-PL,
	// pt-BR, ru-RU, tr-TR, ja-JP, ko-KR, zh-TW, th-TH, en-PH, en-SG
	Locale string
}

func (EsportsClient) getContentStackKeys(params contentstack.Parameters) *contentstack.Keys {
	return contentstack.GetKeys("https://valorantesports.com/news", `a[href^="/news/"]`, &params)
}

func (client EsportsClient) getContentStackParameters(count int) contentstack.Parameters {
	return contentstack.Parameters{
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

	if entry.VideoLink != "" {
		return entry.VideoLink
	}

	return fmt.Sprintf("https://valorantesports.com/%s/%s", utils.TrimSlashes(entry.URL.URL), client.Locale)
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

		results[i] = EsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.BannerSettings.Banner.URL,
			Title:       item.Title,
			URL:         url,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}
