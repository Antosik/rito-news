package riotgames

import (
	"encoding/json"
	"fmt"
	"rito-news/lib/utils"
	"strings"

	"github.com/go-rod/rod"
)

type JobsOfficeEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
	URL  string `json:"url"`
}

type JobsCraftEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
}

type JobsEntry struct {
	Craft    JobsCraftEntry  `json:"craft"`
	Office   JobsOfficeEntry `json:"office"`
	Products string          `json:"products"`
	Title    string          `json:"title"`
	URL      string          `json:"url"`
}

type rawJobsEntry struct {
	Craft    string `json:"craft"`
	CraftID  string `json:"craftId"`
	Office   string `json:"office"`
	OfficeID string `json:"officeId"`
	Products string `json:"products"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type rawJobsResponse struct {
	Jobs []rawJobsEntry `json:"jobs"`
}

type JobsClient struct {
	Locale string
}

func (client JobsClient) loadData() (string, string) {
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

func (client JobsClient) parseData(data string) ([]rawJobsEntry, error) {
	var results rawJobsResponse

	err := json.Unmarshal([]byte(data), &results)
	if err != nil {
		return nil, fmt.Errorf("can't parse data: %w", err)
	}

	return results.Jobs, nil
}

func (client JobsClient) GetItems() ([]JobsEntry, error) {
	data, link := client.loadData()

	items, err := client.parseData(data)
	if err != nil {
		return nil, err
	}

	results := make([]JobsEntry, len(items))
	for i, entry := range items {
		results[i] = JobsEntry{
			Craft: JobsCraftEntry{
				ID:   entry.CraftID,
				Name: entry.Craft,
				More: fmt.Sprintf("%s#craft=%s", link, entry.CraftID),
			},
			Office: JobsOfficeEntry{
				ID:   entry.OfficeID,
				Name: entry.Office,
				More: fmt.Sprintf("%s#office=%s", link, entry.OfficeID),
				URL:  fmt.Sprintf("https://www.riotgames.com/%s/o/%s", client.Locale, entry.OfficeID),
			},
			Products: entry.Products,
			Title:    entry.Title,
			URL:      fmt.Sprintf("https://www.riotgames.com/%s/%s", client.Locale, utils.TrimSlashes(entry.URL)),
		}
	}

	return results, nil
}
