package riotgames

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Antosik/rito-news/internal/utils"
	"github.com/google/uuid"
)

var (
	errNoCareersPageLink     = errors.New("can't find careers page link")
	errNoDataAttributeOnNode = errors.New("can't find data attribute")
)

// Riot Games office entry.
type JobsOfficeEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
	URL  string `json:"url"`
}

// Riot Games craft entry.
type JobsCraftEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	More string `json:"more"`
}

// Riot Games jobs entry.
type JobsEntry struct {
	UID      string          `json:"uid"`
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

// A client that allows to get a list of available vacancies in Riot games.
//
// Source - https://www.riotgames.com/en/work-with-us
type JobsClient struct {
	// Available locales:
	// en, id, ms, pt-br, cs, fr, de, el, hu, it, ja, ko,
	// es-419, pl, ro, ru, zh-cn, es, th, zh-hant, tr, vi
	Locale string
}

func (client JobsClient) loadData() (string, string, error) {
	url := "https://www.riotgames.com/" + client.Locale

	main, err := utils.RunGETHTMLRequest(url)
	if err != nil {
		return "", "", fmt.Errorf("can't get main page html content: %w", err)
	}

	maindoc, err := utils.ReadHTML(main)
	if err != nil {
		return "", "", fmt.Errorf("can't read main page content: %w", err)
	}

	link, linkFound := maindoc.Find(".careers__cta").Attr("href")
	if !linkFound {
		return "", "", errNoCareersPageLink
	}

	if !strings.Contains(link, "https://www.riotgames.com") {
		link = "https://www.riotgames.com/" + utils.TrimSlashes(link)
	}

	jobs, err := utils.RunGETHTMLRequest(link)
	if err != nil {
		return "", "", fmt.Errorf("can't get jobs page html content: %w", err)
	}

	jobsdoc, err := utils.ReadHTML(jobs)
	if err != nil {
		return "", "", fmt.Errorf("can't read jobs content: %w", err)
	}

	data, dataFound := jobsdoc.Find(".js-job-list-wrapper").Attr("data-props")
	if !dataFound {
		return "", "", errNoDataAttributeOnNode
	}

	return data, link, nil
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
	data, link, err := client.loadData()
	if err != nil {
		return nil, err
	}

	items, err := client.parseData(data)
	if err != nil {
		return nil, err
	}

	results := make([]JobsEntry, len(items))

	for i, entry := range items {
		url := fmt.Sprintf("https://www.riotgames.com/%s/%s", client.Locale, utils.TrimSlashes(entry.URL))
		uid := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()

		results[i] = JobsEntry{
			UID: uid,
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
			URL:      url,
		}
	}

	return results, nil
}
