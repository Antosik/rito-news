package val_source

import (
	"encoding/json"
	"errors"
	"fmt"
	"rito-news/sources/base/contentstack"
	utils "rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type VALORANTEsportsEntry struct {
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
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
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

func (client VALORANTEsports) getContentStackItems(count int) ([]VALORANTEsportsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return []VALORANTEsportsEntry{}, err
	}

	items := make([]VALORANTEsportsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return []VALORANTEsportsEntry{}, errors.New("Can't parse item: " + err.Error())
		}
	}

	return items, nil
}

func (client VALORANTEsports) generateNewsLink(entry VALORANTEsportsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	if entry.VideoLink != "" {
		return entry.VideoLink
	}
	return fmt.Sprintf("https://valorantesports.com/%s/%s", utils.TrimSlashes(entry.Url.Url), client.Locale)
}

func (client VALORANTEsports) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.getContentStackItems(count)
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return []abstract.NewsItem{}, errors.New("Can't generate UUID: " + err.Error())
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
