package lor_source

import (
	"fmt"
	"rito-news/sources/base/serverstatus"
	"strings"
)

type LegendsOfRuneterraStatus struct {
	Region string
}

func (client LegendsOfRuneterraStatus) loadItems(locale string) ([]serverstatus.ServerStatusEntry, error) {
	url := fmt.Sprintf(
		"https://bacon.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)
	return serverstatus.GetServerStatusItems(url, locale)
}

func (client LegendsOfRuneterraStatus) generateNewsLink(entry serverstatus.ServerStatusEntry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/lor?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.UID,
	)
}

func (client LegendsOfRuneterraStatus) GetItems(locale string) ([]serverstatus.ServerStatusEntry, error) {
	items, err := client.loadItems(locale)
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].Url = client.generateNewsLink(items[i], locale)
	}

	return items, nil
}
