package lol_source

import (
	"fmt"
	"rito-news/sources/base/serverstatus"
	"rito-news/utils/abstract"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type LeagueOfLegendsStatus struct {
	Region string
}

func (client LeagueOfLegendsStatus) loadItems() (serverstatus.ServerStatusResponse, error) {
	url := fmt.Sprintf(
		"https://lol.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)
	return serverstatus.GetServerStatusItems(url)
}

func (client LeagueOfLegendsStatus) generateNewsLink(entry abstract.NewsItem, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/lol?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.Id,
	)
}

func (client LeagueOfLegendsStatus) GetItems(locale string) ([]abstract.NewsItem, error) {
	status, err := client.loadItems()
	if err != nil {
		return []abstract.NewsItem{}, err
	}

	items := serverstatus.TransformServerStatusToNewsItems(status, locale)
	for i := range items {
		items[i].Url = client.generateNewsLink(items[i], locale)

		id, err := uuid.NewRandomFromReader(strings.NewReader(items[i].Url))
		if err != nil {
			return []abstract.NewsItem{}, fmt.Errorf("can't generate UUID: %w", err)
		}

		items[i].Id = id.String()
	}

	sort.Sort(abstract.ByCreatedAt(items))

	return items, nil
}
