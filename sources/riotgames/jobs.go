package riotgames_source

import (
	"encoding/json"
	"fmt"
	"rito-news/utils"
	"strings"

	"github.com/go-rod/rod"
)

type RiotGamesJobsEntryOffice struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
	Url  string `json:"url"`
}

type RiotGamesJobsEntryCraft struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
}

type RiotGamesJobsEntry struct {
	Craft    RiotGamesJobsEntryCraft  `json:"craft"`
	Office   RiotGamesJobsEntryOffice `json:"office"`
	Products string                   `json:"products"`
	Title    string                   `json:"title"`
	Url      string                   `json:"url"`
}

type riotGamesJobsResponseAPIResponseEntry struct {
	Craft    string `json:"craft"`
	CraftId  string `json:"craftId"`
	Office   string `json:"office"`
	OfficeId string `json:"officeId"`
	Products string `json:"products"`
	Title    string `json:"title"`
	Url      string `json:"url"`
}

type riotGamesJobsResponseAPIResponse struct {
	Jobs []riotGamesJobsResponseAPIResponseEntry `json:"jobs"`
}

type RiotGamesJobs struct {
	Locale string
}

func (client RiotGamesJobs) loadData() (string, string) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(fmt.Sprintf("https://www.riotgames.com/%s", client.Locale))
	defer page.MustClose()

	link := *page.MustElement(".careers__cta").MustAttribute("href")
	if !strings.Contains(link, "https://www.riotgames.com") {
		link = "https://www.riotgames.com" + utils.TrimSlashes(link)
	}
	page.MustNavigate(link)

	data := page.MustElement(".js-job-list-wrapper").MustAttribute("data-props")

	return *data, link
}

func (client RiotGamesJobs) parseData(data string) ([]riotGamesJobsResponseAPIResponseEntry, error) {
	var results riotGamesJobsResponseAPIResponse

	err := json.Unmarshal([]byte(data), &results)
	if err != nil {
		return nil, fmt.Errorf("can't parse data: %w", err)
	}

	return results.Jobs, nil
}

func (client RiotGamesJobs) GetItems() ([]RiotGamesJobsEntry, error) {
	data, link := client.loadData()

	items, err := client.parseData(data)
	if err != nil {
		return nil, err
	}

	results := make([]RiotGamesJobsEntry, len(items))
	for i, entry := range items {
		results[i] = RiotGamesJobsEntry{
			Craft: RiotGamesJobsEntryCraft{
				Id:   entry.CraftId,
				Name: entry.Craft,
				More: fmt.Sprintf("%s#craft=%s", link, entry.CraftId),
			},
			Office: RiotGamesJobsEntryOffice{
				Id:   entry.OfficeId,
				Name: entry.Office,
				More: fmt.Sprintf("%s#office=%s", link, entry.OfficeId),
				Url:  fmt.Sprintf("https://www.riotgames.com/%s/o/%s", client.Locale, entry.OfficeId),
			},
			Products: entry.Products,
			Title:    entry.Title,
			Url:      fmt.Sprintf("https://www.riotgames.com/%s/%s", client.Locale, utils.TrimSlashes(entry.Url)),
		}
	}

	return results, nil
}
