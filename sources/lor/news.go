package lor_source

import (
	"encoding/json"
	"errors"
	"fmt"
	"rito-news/sources/base/contentstack"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LegendsOfRuneterraNewsEntry struct {
	UID         string `json:"uid"`
	ArticleTags []struct {
		MachineName string `json:"machine_name"`
		Title       string `json:"title"`
	} `json:"article_tags"`
	Categories []struct {
		MachineName string `json:"machine_name"`
		Title       string `json:"title"`
	} `json:"category"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Author  []struct {
		Title string `json:"title"`
	} `json:"author"`
	Image struct {
		Url string `json:"url"`
	} `json:"cover_image"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
}

type LegendsOfRuneterraNews struct {
	Locale string
}

func (LegendsOfRuneterraNews) getContentStackKeys(params contentstack.ContentStackQueryParameters) *contentstack.ContentStackKeys {
	return contentstack.GetContentStackKeys("https://playruneterra.com/en-us/news/", ".page ul li", &params)
}

func (client LegendsOfRuneterraNews) getContentStackParameters(count int) contentstack.ContentStackQueryParameters {
	return contentstack.ContentStackQueryParameters{
		ContentType: "news_2",
		Locale:      client.Locale,
		Count:       count,
		Environment: "live",
		Filters: map[string][]string{
			"query": {`{"hide_from_newsfeeds":{"$ne":true}}`},
			"only[BASE][]": {
				"title",
				"_content_type_uid",
				"cover_image",
				"author",
				"article_tags",
				"category",
				"date",
				"summary",
				"external_link",
				"url",
			},
			"include[]": {
				"article_tags",
				"author",
				"category",
			},
			"only[article_tags][]": {
				"machine_name",
				"title",
				"url",
			},
			"only[author][]": {
				"title",
				"image",
			},
			"only[category][]": {
				"machine_name",
				"title",
				"url",
			},
		},
	}
}

func (client LegendsOfRuneterraNews) getContentStackItems(count int) ([]LegendsOfRuneterraNewsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetContentStackItems(keys, &params)
	if err != nil {
		return []LegendsOfRuneterraNewsEntry{}, err
	}

	items := make([]LegendsOfRuneterraNewsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return []LegendsOfRuneterraNewsEntry{}, errors.New("Can't parse item: " + err.Error())
		}
	}

	return items, nil
}

func (client LegendsOfRuneterraNews) generateNewsLink(entry LegendsOfRuneterraNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://playruneterra.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client LegendsOfRuneterraNews) GetItems(count int) ([]abstract.NewsItem, error) {
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

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Summary,
			Url:       url,
			Author:    strings.Join(authors, ","),
			Image:     item.Image.Url,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
