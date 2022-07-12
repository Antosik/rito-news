package riotgames

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Antosik/rito-news/internal/utils"

	"github.com/PuerkitoBio/goquery"
)

type NewsEntry struct {
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

type NewsClient struct {
	Locale string
}

func (client NewsClient) initialLoad() ([]string, string, error) {
	main, err := utils.RunGETHTMLRequest(fmt.Sprintf("https://www.riotgames.com/%s", client.Locale))
	if err != nil {
		return nil, "", err
	}

	maindoc, err := utils.ReadHTML(main)
	if err != nil {
		return nil, "", fmt.Errorf("can't read main content: %w", err)
	}

	link, linkFound := maindoc.Find(".whats-happening__cta").Attr("href")
	if !linkFound {
		return nil, "", fmt.Errorf("can't find careers page link")
	}

	if !strings.Contains(link, "https://www.riotgames.com") {
		link = "https://www.riotgames.com/" + utils.TrimSlashes(link)
	}

	news, err := utils.RunGETHTMLRequest(link)
	if err != nil {
		return nil, "", err
	}

	newsdoc, err := utils.ReadHTML(news)
	if err != nil {
		return nil, "", fmt.Errorf("can't read news content: %w", err)
	}

	ids, idsFound := newsdoc.Find(".js-load-more").Attr("data-load-more-ids")
	if !idsFound {
		return nil, "", fmt.Errorf("can't find ids to load: %w", err)
	}

	newsElements := newsdoc.Find(".js-explore-hero-wrapper .content-center, .widget__wrapper--maxigrid .grid__item")
	newsHTML := make([]string, len(newsElements.Nodes))

	for i := range newsElements.Nodes {
		el := newsElements.Eq(i)
		newsHTML[i], _ = el.Html()
	}

	return strings.Split(ids, ","), strings.Join(newsHTML, ""), nil
}

func (client NewsClient) loadNewsWithIds(ids []string) (string, error) {
	widget := fmt.Sprintf(`{"loadMorePageSize":%d,"loadMoreMethod":"button"}`, len(ids))
	url := fmt.Sprintf(
		`https://www.riotgames.com/%s/api/load-more/maxi-grid?ids=%s&widget=%s`,
		client.Locale,
		strings.Join(ids, ","),
		widget,
	)

	body, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(body, `\"`, `"`), nil
}

func (NewsClient) extractNewsFromHTML(html string) ([]NewsEntry, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("can't read news content: %w", err)
	}

	items := doc.Find(".summary")
	news := make([]NewsEntry, items.Size())

	for i := range items.Nodes {
		el := items.Eq(i)

		dateStr := el.Find(".summary__date").Text()

		var (
			date time.Time
			err  error
		)

		date, err = utils.ParseDateTimeWithLayouts(dateStr, []string{
			"02/01/2006",
			"Jan 2, 2006",
			"January 2, 2006",
			"2006/01/02",
			"02.01.2006",
			"2 January, 2006",
		})
		if err != nil {
			return nil, fmt.Errorf("can't parse article date: %w", err)
		}

		url, _ := el.Find("a").Attr("href")
		if !strings.Contains(url, "https://www.riotgames.com") {
			url = "https://www.riotgames.com/" + utils.TrimSlashes(url)
		}

		news[i] = NewsEntry{
			Category:    el.Find(".eyebrow span").Text(),
			Date:        date,
			Description: el.Find(".summary__sell").Text(),
			Image:       el.Find("img").AttrOr("src", ""),
			Title:       el.Find("h3 span").Text(),
			URL:         url,
		}
	}

	return news, nil
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	ids, initialsNews, err := client.initialLoad()
	if err != nil {
		return nil, err
	}

	items, err := client.extractNewsFromHTML(initialsNews)
	if err != nil {
		return nil, err
	}

	if count > len(items) {
		idsToLoadCount := utils.MinInt(count-len(items), len(ids))

		news, err := client.loadNewsWithIds(ids[:idsToLoadCount])
		if err != nil {
			return nil, err
		}

		additionalNews, err := client.extractNewsFromHTML(news)
		if err != nil {
			return nil, err
		}

		items = append(items, additionalNews...)
	} else {
		items = items[:count]
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	return items, nil
}
