package lol

import (
	"encoding/json"
	"fmt"
	"rito-news/lib/contentstack"
	"rito-news/lib/utils"
	"time"
)

type LeagueOfLegendsEsportsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type leagueOfLegendsEsportsAPIResponseEntry struct {
	UID    string `json:"uid"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	HeaderImage  struct {
		Url string `json:"url"`
	} `json:"header_image"`
	Intro string `json:"intro"`
	Title string `json:"title"`
	Url   string `json:"url"`
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

func (client LeagueOfLegendsEsports) getContentStackItems(count int) ([]leagueOfLegendsEsportsAPIResponseEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]leagueOfLegendsEsportsAPIResponseEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (LeagueOfLegendsEsports) generateNewsLink(entry leagueOfLegendsEsportsAPIResponseEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://lolesports.com/%s", utils.TrimSlashes(entry.Url))
}

func (client LeagueOfLegendsEsports) GetItems(count int) ([]LeagueOfLegendsEsportsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]LeagueOfLegendsEsportsEntry, len(items))
	for i, item := range items {
		url := client.generateNewsLink(item)

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		results[i] = LeagueOfLegendsEsportsEntry{
			UID:         item.UID,
			Authors:     authors,
			Date:        item.Date,
			Description: item.Intro,
			Image:       item.HeaderImage.Url,
			Title:       item.Title,
			Url:         url,
		}
	}

	return results, nil
}
