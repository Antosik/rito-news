package lor

import (
	"encoding/json"
	"fmt"
	"rito-news/internal/contentstack"
	"rito-news/internal/utils"
	"time"
)

type NewsEntry struct {
	UID         string    `json:"uid"`
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type rawNewsEntry struct {
	UID         string `json:"uid"`
	ArticleTags []struct {
		Title string `json:"title"`
	} `json:"article_tags"`
	Author []struct {
		Title string `json:"title"`
	} `json:"author"`
	Category []struct {
		Title string `json:"title"`
	} `json:"category"`
	CoverImage struct {
		URL string `json:"url"`
	} `json:"cover_image"`
	Date         time.Time `json:"date"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	URL          struct {
		URL string `json:"url"`
	} `json:"url"`
}

type NewsClient struct {
	Locale string
}

func (NewsClient) getContentStackKeys(params contentstack.Parameters) *contentstack.Keys {
	return contentstack.GetKeys("https://playruneterra.com/en-us/news/", ".page ul li", &params)
}

func (client NewsClient) getContentStackParameters(count int) contentstack.Parameters {
	return contentstack.Parameters{
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

func (client NewsClient) getContentStackItems(count int) ([]rawNewsEntry, error) {
	params := client.getContentStackParameters(count)
	keys := client.getContentStackKeys(params)

	rawitems, err := contentstack.GetItems(keys, &params)
	if err != nil {
		return nil, err
	}

	items := make([]rawNewsEntry, len(rawitems))

	for i, raw := range rawitems {
		err := json.Unmarshal(raw, &items[i])
		if err != nil {
			return nil, fmt.Errorf("can't parse item: %w", err)
		}
	}

	return items, nil
}

func (client NewsClient) getLinkForEntry(entry rawNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}

	return fmt.Sprintf("https://playruneterra.com/%s/%s/", client.Locale, utils.TrimSlashes(entry.URL.URL))
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	items, err := client.getContentStackItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]NewsEntry, len(items))

	for i, item := range items {
		url := client.getLinkForEntry(item)

		authors := make([]string, len(item.Author))
		for i, author := range item.Author {
			authors[i] = author.Title
		}

		categories := make([]string, len(item.Category))
		for i, category := range item.Category {
			categories[i] = category.Title
		}

		tags := make([]string, len(item.ArticleTags))
		for i, tag := range item.ArticleTags {
			tags[i] = tag.Title
		}

		results[i] = NewsEntry{
			UID:         item.UID,
			Authors:     authors,
			Categories:  categories,
			Date:        item.Date,
			Description: item.Summary,
			Image:       item.CoverImage.URL,
			Tags:        tags,
			Title:       item.Title,
			URL:         url,
		}
	}

	return results, nil
}
