package riotgames

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

var errNoNewsPageLink = errors.New("can't find news page link")

// RiotGames news entry.
type NewsEntry struct {
	UID         string    `json:"uid"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
}

// A client that allows to get official RiotGames news.
//
// Source - https://www.riotgames.com/en/news
type NewsClient struct {
	// Available locales:
	// en, id, ms, pt-br, cs, fr, de, el, hu, it, ja, ko,
	// es-419, pl, ro, ru, zh-cn, es, th, zh-hant, tr, vi
	Locale string
}

func (client NewsClient) initialLoad() ([]string, string, error) {
	url := "https://www.riotgames.com/" + client.Locale

	main, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return nil, "", fmt.Errorf("can't get main page html content: %w", err)
	}

	maindoc, err := utils.ReadHTML(main)
	if err != nil {
		return nil, "", fmt.Errorf("can't read main page content: %w", err)
	}

	link, linkFound := maindoc.Find(".whats-happening__cta").Attr("href")
	if !linkFound {
		return nil, "", errNoNewsPageLink
	}

	if !strings.Contains(link, "https://www.riotgames.com") {
		link = "https://www.riotgames.com/" + utils.TrimSlashes(link)
	}

	news, err := utils.RunGETHTMLRequest(link)
	if err != nil {
		return nil, "", fmt.Errorf("can't get news page html content: %w", err)
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

func (client NewsClient) loadNewsWithIDs(ids []string) (string, error) {
	widget := fmt.Sprintf(`{"loadMorePageSize":%d,"loadMoreMethod":"button"}`, len(ids))
	url := fmt.Sprintf(
		`https://www.riotgames.com/%s/api/load-more/maxi-grid?ids=%s&widget=%s`,
		client.Locale,
		strings.Join(ids, ","),
		widget,
	)

	body, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return "", fmt.Errorf("can't load more news html: %w", err)
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

	for index := range items.Nodes {
		node := items.Eq(index)

		dateStr := node.Find(".summary__date").Text()

		var (
			date time.Time
			err  error
		)

		date, err = utils.ParseDateTimeWithLayouts(dateStr, []string{
			"02/01/2006",
			"02/01/200602/01/2006",
			"Jan 2, 2006",
			"Jan 2, 2006Jan 2, 2006",
			"January 2, 2006",
			"January 2, 2006January 2, 2006",
			"2006/01/02",
			"2006/01/022006/01/02",
			"02.01.2006",
			"02.01.200602.01.2006",
			"2 January, 2006",
		})
		if err != nil {
			return nil, fmt.Errorf("can't parse article date: %w", err)
		}

		url, _ := node.Find("a").Attr("href")
		if !strings.Contains(url, "https://www.riotgames.com") {
			url = "https://www.riotgames.com/" + utils.TrimSlashes(url)
		}

		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

		news[index] = NewsEntry{
			UID:         uid,
			Category:    node.Find(".eyebrow span").Text(),
			Date:        date,
			Description: node.Find(".summary__sell").Text(),
			Image:       node.Find("img").AttrOr("src", ""),
			Title:       node.Find("h3 span").Text(),
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

		news, err := client.loadNewsWithIDs(ids[:idsToLoadCount])
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
