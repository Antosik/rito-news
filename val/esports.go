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
	Tags        []string  `json:"tags"`
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
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Event       []struct {
		Title string `json:"title"`
	} `json:"event"`
	ExternalLink string `json:"external_link"`
	Title        string `json:"title"`
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
	// en-us, en-gb, en-au, de-de, es-es, es-mx, fr-fr, it-it, pl-pl,
	// pt-br, ru-ru, tr-tr, ja-jp, ko-kr, zh-tw, th-th, en-ph, en-sg
	// id-id, vi-vn
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
				"event",
				"external_link",
				"url",
				"video_link",
			},
			"include[]": {
				"authors",
				"event",
			},
			"only[authors][]": {
				"title",
			},
			"only[event][]": {
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

		tags := make([]string, len(item.Event))
		for i, event := range item.Event {
			tags[i] = event.Title
		}

		results[i] = EsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Date:        item.Date,
			Description: item.Description,
			Image:       item.BannerSettings.Banner.URL,
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
