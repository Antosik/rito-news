package lol

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/contentstack"
	"github.com/Antosik/rito-news/internal/utils"
)

// League of Legends esports news entry
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
	UID    string `json:"uid"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	HeaderImage  struct {
		URL string `json:"url"`
	} `json:"header_image"`
	Intro  string `json:"intro"`
	League []struct {
		Title string `json:"title"`
	} `json:"league"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// A client that allows to get League of Legends esports news.
//
// Source - https://lolesports.com/news
type EsportsClient struct {
	// Available locales:
	// en-us, en-gb, de-de, es-es, es-mx, fr-fr, it-it, pl-pl, pt-br,
	// ru-ru, tr-tr, ja-jp, ko-kr, zh-tw, th-th, en-ph, en-sg
	Locale string
}

func (EsportsClient) getContentStackKeys(params contentstack.Parameters) *contentstack.Keys {
	return contentstack.GetKeys("https://lolesports.com/news", ".News .content-block", &params)
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
				"header_image",
				"author",
				"league",
				"date",
				"intro",
				"external_link",
				"url",
			},
			"include[]": {
				"author",
				"league",
			},
			"only[author][]": {
				"title",
			},
			"only[league][]": {
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

func (EsportsClient) getLinkForEntry(entry rawEsportsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	return fmt.Sprintf("https://lolesports.com/%s", utils.TrimSlashes(entry.URL))
}

func (client EsportsClient) GetItems(count int) ([]EsportsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]EsportsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		tags := make([]string, len(item.League))
		for i, league := range item.League {
			tags[i] = league.Title
		}

		results[i] = EsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Date:        item.Date,
			Description: item.Intro,
			Image:       item.HeaderImage.URL,
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
