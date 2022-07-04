package lol

import (
	"fmt"
	"rito-news/lib/serverstatus"
	"strings"
)

type LeagueOfLegendsStatus struct {
	Region string
}

func (client LeagueOfLegendsStatus) loadItems(locale string) ([]serverstatus.ServerStatusEntry, error) {
	url := fmt.Sprintf(
		"https://lol.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)
	return serverstatus.GetServerStatusItems(url, locale)
}

func (client LeagueOfLegendsStatus) generateNewsLink(entry serverstatus.ServerStatusEntry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/lol?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.UID,
	)
}

func (client LeagueOfLegendsStatus) GetItems(locale string) ([]serverstatus.ServerStatusEntry, error) {
	items, err := client.loadItems(locale)
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].Url = client.generateNewsLink(items[i], locale)
	}

	return items, nil
}
