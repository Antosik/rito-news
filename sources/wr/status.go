package wr_source

import (
	"fmt"
	"rito-news/sources/base/serverstatus"
	"rito-news/utils/abstract"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type WildRiftStatus struct {
	Region string
}

func (client WildRiftStatus) loadItems() (serverstatus.ServerStatusResponse, error) {
	url := fmt.Sprintf(
		"https://wildrift.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)
	return serverstatus.GetServerStatusItems(url)
}

func (client WildRiftStatus) generateNewsLink(entry abstract.NewsItem, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/wildrift?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.Id,
	)
}

func (client WildRiftStatus) GetItems(locale string) ([]abstract.NewsItem, error) {
	status, err := client.loadItems()
	if err != nil {
		return nil, err
	}

	items := serverstatus.TransformServerStatusToNewsItems(status, locale)
	for i := range items {
		items[i].Url = client.generateNewsLink(items[i], locale)

		id, err := uuid.NewRandomFromReader(strings.NewReader(items[i].Url))
		if err != nil {
			return nil, fmt.Errorf("can't generate UUID: %w", err)
		}

		items[i].Id = id.String()
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
