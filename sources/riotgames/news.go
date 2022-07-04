package riotgames_source

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"rito-news/utils"
	"rito-news/utils/abstract"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/google/uuid"
)

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
		return "", errors.New("Can't load more news: " + err.Error())
	}

	body, _ := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Can't parse news body: " + err.Error())
	}

	return strings.ReplaceAll(string(body), `\"`, `"`), nil
}

func (RiotGamesNews) extractNewsFromHTML(html string) ([]abstract.NewsItem, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return []abstract.NewsItem{}, errors.New("Can't read news content: " + err.Error())
	}

	items := doc.Find(".summary")
	news := make([]abstract.NewsItem, items.Size())

	for i := range items.Nodes {
		el := items.Eq(i)

		title := el.Find("h3 span").Text()
		summary := el.Find(".summary__sell").Text()
		category := el.Find(".eyebrow span").Text()
		image := el.Find("img").AttrOr("src", "")

		dateStr := el.Find(".summary__date").Text()
		var (
			date time.Time
			err  error
		)
		date, err = time.Parse("02/01/2006", dateStr)
		if err != nil {
			date, err = time.Parse("Jan 2, 2006", dateStr)
			if err != nil {
				return []abstract.NewsItem{}, errors.New("Can't parse article date: " + err.Error())
			}
		}

		url, _ := el.Find("a").Attr("href")
		if !strings.Contains(url, "https://www.riotgames.com") {
			url = "https://www.riotgames.com/" + utils.TrimSlashes(url)
		}

		id, err := uuid.NewRandomFromReader(strings.NewReader(url))
		if err != nil {
			return []abstract.NewsItem{}, errors.New("Can't generate UUID: " + err.Error())
		}

		news[i] = abstract.NewsItem{
			Id:        id.String(),
			Url:       url,
			Title:     title,
			Summary:   summary,
			Author:    "Riot Games",
			Category:  category,
			Image:     image,
			CreatedAt: date,
			UpdatedAt: date,
		}
	}

	return news, nil
}

func (client RiotGamesNews) GetItems(count int) ([]abstract.NewsItem, error) {
	ids, initialsNews := client.initialLoad()

	items, err := client.extractNewsFromHTML(initialsNews)
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	if count > len(items) {
		idsToLoadCount := count - len(items)
		if len(ids) < idsToLoadCount {
			idsToLoadCount = len(ids)
		}

		news, err := client.loadNewsWithIds(ids[:idsToLoadCount])
		if err != nil {
			return []abstract.NewsItem{}, err
		}

		additionalNews, err := client.extractNewsFromHTML(news)
		if err != nil {
			return []abstract.NewsItem{}, err
		}

		items = append(items, additionalNews...)
	} else {
		items = items[:count]
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
