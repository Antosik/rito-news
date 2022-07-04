package val_source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type VALORANTNewsEntry struct {
	Banner struct {
		Url string `json:"url"`
	} `json:"banner"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	ExternalLink string    `json:"external_link"`
	Title        string    `json:"title"`
	UID          string    `json:"uid"`
	Url          struct {
		Url string `json:"url"`
	} `json:"url"`
}

type VALORANTNewsResponse struct {
	Result struct {
		Data struct {
			AllContentstackArticles struct {
				Nodes []VALORANTNewsEntry `json:"nodes"`
			} `json:"allContentstackArticles"`
		} `json:"data"`
	} `json:"result"`
}

type VALORANTNews struct {
	Locale string
}

func (client VALORANTNews) loadItems(count int) ([]VALORANTNewsEntry, error) {
	url := fmt.Sprintf(
		"https://playvalorant.com/page-data/%s/news/page-data.json",
		client.Locale,
	)
	res, err := http.Get(url)
	if err != nil {
		return []VALORANTNewsEntry{}, fmt.Errorf("can't load news: %w", err)
	}
	defer res.Body.Close()

	var response VALORANTNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []VALORANTNewsEntry{}, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Result.Data.AllContentstackArticles.Nodes[:count], nil
}

func (client VALORANTNews) generateNewsLink(entry VALORANTNewsEntry) string {
	if entry.ExternalLink != "" {
		return entry.ExternalLink
	}
	return fmt.Sprintf("https://playvalorant.com/%s/%s", client.Locale, utils.TrimSlashes(entry.Url.Url))
}

func (client VALORANTNews) GetItems(count int) ([]abstract.NewsItem, error) {
	stackItems, err := client.loadItems(count)
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	items := make([]abstract.NewsItem, len(stackItems))

	for i, item := range stackItems {
		url := client.generateNewsLink(item)

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return []abstract.NewsItem{}, fmt.Errorf("can't generate UUID: %w", err)
		}

		items[i] = abstract.NewsItem{
			Id:        id.String(),
			Title:     item.Title,
			Summary:   item.Description,
			Url:       url,
			Image:     item.Banner.Url,
			CreatedAt: item.Date,
			UpdatedAt: item.Date,
		}
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
