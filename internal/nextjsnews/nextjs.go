// Package nextjs implements the way to pull the news from nextjs-based news pages
package nextjsnews

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

type Item struct {
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

type rawItem struct {
	Title  string    `json:"title"`
	Date   time.Time `json:"publishedAt"`
	Action struct {
		Type    string `json:"type"` // 'weblink', 'youtube_video'
		Payload struct {
			URL string `json:"url"`
		} `json:"payload"`
	} `json:"action"`
	Media struct {
		Type string `json:"type"` // 'image'
		URL  string `json:"url"`
	} `json:"media"`
	Description struct {
		Type string `json:"type"` // 'html'
		Body string `json:"body"`
	} `json:"description"`
	Category struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		MachineName string `json:"machineName"`
	} `json:"category"`
}

type rawDataBladeResponse struct {
	Type       string    `json:"type"`
	FragmentID string    `json:"fragmentId"` // should be 'news'
	Items      []rawItem `json:"items"`
}

type rawDataResponse struct {
	Props struct {
		PageProps struct {
			Page struct {
				Blades []rawDataBladeResponse `json:"blades"`
			} `json:"page"`
		} `json:"pageProps"`
	} `json:"props"`
}

type Parser struct {
	BaseURL string
	Locale  string
}

func (parser Parser) loadData() ([]rawItem, error) {
	url := fmt.Sprintf(
		"%s/%s/news/",
		parser.BaseURL,
		parser.Locale,
	)

	body, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return nil, fmt.Errorf("can't get html content: %w", err)
	}

	doc, err := utils.ReadHTML(body)
	if err != nil {
		return nil, fmt.Errorf("can't read news page content: %w", err)
	}

	content := doc.Find("#__NEXT_DATA__").Text()

	var response rawDataResponse
	if err := json.NewDecoder(strings.NewReader(content)).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	for _, item := range response.Props.PageProps.Page.Blades {
		if item.FragmentID != "news" || len(item.Items) == 0 {
			continue
		}

		return item.Items, nil
	}

	return nil, fmt.Errorf("can't find news data: %w", err)
}

func (parser Parser) getLinkForItem(item rawItem) string {
	link := item.Action.Payload.URL

	switch linkType := item.Action.Type; linkType {
	case "weblink":
		if strings.HasPrefix(link, "http") {
			return link
		}

		return fmt.Sprintf("%s/%s/", parser.BaseURL, utils.TrimSlashes(item.Action.Payload.URL))
	case "youtube_video":
		return link
	default:
		return fmt.Sprintf("%s/%s/news", parser.BaseURL, parser.Locale)
	}
}

func (parser Parser) transformRawItems(rawItems []rawItem) []Item {
	items := make([]Item, len(rawItems))

	for index, rawItem := range rawItems {
		url := parser.getLinkForItem(rawItem)
		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

		authors := make([]string, 0)
		categories := []string{rawItem.Category.Title}
		tags := make([]string, 0)

		items[index] = Item{
			UID:         uid,
			Authors:     authors,
			Categories:  categories,
			Date:        rawItem.Date,
			Description: rawItem.Description.Body,
			Image:       rawItem.Media.URL,
			Tags:        tags,
			Title:       rawItem.Title,
			URL:         url,
		}
	}

	return items
}

func (parser Parser) GetItems(count int) ([]Item, error) {
	// Load raw items
	rawItems, err := parser.loadData()
	if err != nil {
		return nil, fmt.Errorf("can't load news page data: %w", err)
	}

	// Sort in case of same publish date
	sort.Slice(rawItems, func(i, j int) bool {
		datecomp := rawItems[i].Date.Compare(rawItems[j].Date)

		if datecomp == 0 {
			return rawItems[i].Title < rawItems[j].Title
		}

		return datecomp > 0
	})

	// Calculate max items count
	sliceSize := utils.MinInt(count, len(rawItems))

	return parser.transformRawItems(rawItems[:sliceSize]), nil
}
