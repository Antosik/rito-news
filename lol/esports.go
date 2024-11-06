package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

// League of Legends esports news entry.
type EsportsEntry struct {
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

type rawEsportsEntry struct {
	ID      string `json:"_id"`
	Authors []struct {
		Title string `json:"externalTitle"`
	} `json:"authors"`
	BannerMedia struct {
		SanityImage struct {
			Asset struct {
				URL string `json:"url"`
			} `json:"asset"`
		} `json:"sanityImage"`
	} `json:"bannerMedia"`
	Description   string `json:"description"`
	ExternalTitle string `json:"externalTitle"`
	ExternalURL   string `json:"externalUrl"`
	Path          struct {
		Current string `json:"current"`
	} `json:"path"`
	PublishingDates struct {
		DisplayedPublishDate time.Time `json:"displayedPublishDate"`
	} `json:"publishingDates"`
}

type rawResponse struct {
	Data struct {
		AllArticle []rawEsportsEntry `json:"allArticle"`
	} `json:"data"`
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

func (client EsportsClient) loadItems(count int) ([]rawEsportsEntry, error) {
	operationName := "LoadMoreNewsList"
	//nolint:lll
	variables := url.QueryEscape(fmt.Sprintf(`{"limit":%d,"offset":0,"sort":[{"publishingDates":{"displayedPublishDate":"DESC"}}],"where":{"channel":{"_ref":{"eq":"channel.league_of_legends_esports_website.%s"}}}}`, count, client.Locale))
	//nolint:lll
	extensions := url.QueryEscape(`{"persistedQuery":{"version":1,"sha256Hash":"790cc5ed50dd93011f92b3ed1bfcb98c70b5353f5fb718e90590e08a2e9124ff"}}`)
	query := "operationName=" + operationName + "&" + "variables=" + variables + "&" + "extensions=" + extensions

	//exhaustruct:ignore
	url := url.URL{
		Scheme:   "https",
		Host:     "lolesports.com",
		Path:     "/api/gql",
		RawQuery: query,
	}

	req, err := utils.NewGETJSONRequest(url.String())
	if err != nil {
		return nil, fmt.Errorf("can't get json content: %w", err)
	}

	req.Header.Set("Apollographql-Client-Name", "Esports Web")
	req.Header.Set("Apollographql-Client-Version", "0c4923c")
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{} //nolint:exhaustruct

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful request: %w", err)
	}
	defer resp.Body.Close()

	var response rawResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	return response.Data.AllArticle, nil
}

func (EsportsClient) getLinkForEntry(entry rawEsportsEntry) string {
	if entry.ExternalURL != "" {
		return entry.ExternalURL
	}

	return "https://lolesports.com/%s" + utils.TrimSlashes(entry.Path.Current)
}

func (client EsportsClient) GetItems(count int) ([]EsportsEntry, error) {
	items, err := client.loadItems(count)
	if err != nil {
		return nil, err
	}

	results := make([]EsportsEntry, len(items))

	for index, item := range items {
		url := client.getLinkForEntry(item)
		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

		authors := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authors[i] = author.Title
		}

		categories := make([]string, 0)
		tags := make([]string, 0)

		results[index] = EsportsEntry{
			UID:         uid,
			Authors:     authors,
			Categories:  categories,
			Date:        item.PublishingDates.DisplayedPublishDate,
			Description: item.Description,
			Image:       item.BannerMedia.SanityImage.Asset.URL,
			Tags:        tags,
			Title:       item.ExternalTitle,
			URL:         url,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}
