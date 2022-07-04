package lol_source

import (
	"encoding/json"
	"fmt"
	"rito-news/sources/base/contentstack"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LeagueOfLegendsEsportsEntry struct {
	UID    string `json:"uid"`
	Title  string `json:"title"`
	Intro  string `json:"intro"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Image struct {
		Url string `json:"url"`
	} `json:"header_image"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	Url          string    `json:"url"`
}

type LeagueOfLegendsEsports struct {
	Locale string
}

func (LeagueOfLegendsEsports) getContentStackKeys(params contentstack.ContentStackQueryParameters) *contentstack.ContentStackKeys {
	return contentstack.GetContentStackKeys("https://lolesports.com/news", ".News .content-block", &params)
}

func (client LeagueOfLegendsEsports) getContentStackParameters(count int) contentstack.ContentStackQueryParameters {
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
				"header_image",
				"author",
				"date",
				"intro",
				"external_link",
				"url",
			},
			"include[]": {
				"author",
			},
			"only[author][]": {
				"title",
			},
		},
	}
}

func (client LeagueOfLegendsEsports) getContentStackItems(count int) ([]LeagueOfLegendsEsportsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]LeagueOfLegendsEsportsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (LeagueOfLegendsEsports) generateNewsLink(entry LeagueOfLegendsEsportsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://lolesports.com/%s", utils.TrimSlashes(entry.Url))
}

func (client LeagueOfLegendsEsports) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return nil, fmt.Errorf("can't generate UUID: %w", err)
		}

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Intro,
			Url:       url,
			Author:    strings.Join(authors, ","),
			Image:     item.Image.Url,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	return items, nil
}
