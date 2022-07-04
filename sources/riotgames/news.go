package riotgames_source

import (
	"fmt"
	"io"
	"net/http"
	"rito-news/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

type RiotGamesNewsEntry struct {
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
}

type RiotGamesNews struct {
	Locale string
}

func (client RiotGamesNews) initialLoad() ([]string, string) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(fmt.Sprintf("https://www.riotgames.com/%s", client.Locale))
	defer page.MustClose()

	link := *page.MustElement(".whats-happening__cta").MustAttribute("href")
	if !strings.Contains(link, "https://www.riotgames.com") {
		link = "https://www.riotgames.com" + link
	}
	page.MustNavigate(link)

	ids := page.MustElement(".js-load-more").MustAttribute("data-load-more-ids")

	news := page.MustElements(".js-explore-hero-wrapper .content-center, .widget__wrapper--maxigrid .grid__item")
	newsHTML := make([]string, len(news))
	for i, newsItem := range news {
		newsHTML[i], _ = newsItem.HTML()
	}

	return strings.Split(*ids, ","), strings.Join(newsHTML, "")
}

func (client RiotGamesNews) loadNewsWithIds(ids []string) (string, error) {
	url := fmt.Sprintf(
		`https://www.riotgames.com/%s/api/load-more/maxi-grid?ids=%s&widget={"loadMorePageSize":%d,"loadMoreMethod":"button"}`,
		client.Locale,
		strings.Join(ids, ","),
		len(ids),
	)

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("can't load more news: %w", err)
	}

	body, _ := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("can't parse news body: %w", err)
	}

	return strings.ReplaceAll(string(body), `\"`, `"`), nil
}

func (RiotGamesNews) extractNewsFromHTML(html string) ([]RiotGamesNewsEntry, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("can't read news content: %w", err)
	}

	items := doc.Find(".summary")
	news := make([]RiotGamesNewsEntry, items.Size())

	for i := range items.Nodes {
		el := items.Eq(i)

		dateStr := el.Find(".summary__date").Text()

		var (
			date time.Time
			err  error
		)
		date, err = time.Parse("02/01/2006", dateStr)
		if err != nil {
			date, err = time.Parse("Jan 2, 2006", dateStr)
			if err != nil {
				return nil, fmt.Errorf("can't parse article date: %w", err)
			}
		}

		url, _ := el.Find("a").Attr("href")
		if !strings.Contains(url, "https://www.riotgames.com") {
			url = "https://www.riotgames.com/" + utils.TrimSlashes(url)
		}

		news[i] = RiotGamesNewsEntry{
			Category:    el.Find(".eyebrow span").Text(),
			Date:        date,
			Description: el.Find(".summary__sell").Text(),
			Image:       el.Find("img").AttrOr("src", ""),
			Title:       el.Find("h3 span").Text(),
			Url:         url,
		}
	}

	return news, nil
}

func (client RiotGamesNews) GetItems(count int) ([]RiotGamesNewsEntry, error) {
	ids, initialsNews := client.initialLoad()

	items, err := client.extractNewsFromHTML(initialsNews)
	if err != nil {
		return nil, err
	}

	if count > len(items) {
		idsToLoadCount := count - len(items)
		if len(ids) < idsToLoadCount {
			idsToLoadCount = len(ids)
		}

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

	return items, nil
}
