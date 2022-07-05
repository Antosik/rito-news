package val

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Antosik/rito-news/internal/serverstatus"
)

type StatusClient struct {
	Region string
}

func (client StatusClient) loadItems(locale string) ([]serverstatus.Entry, error) {
	url := fmt.Sprintf(
		"https://valorant.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)

	return serverstatus.GetItems(url, locale)
}

func (client StatusClient) getLinkForEntry(entry serverstatus.Entry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/valorant?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.UID,
	)
}

func (client StatusClient) GetItems(locale string) ([]serverstatus.Entry, error) {
	items, err := client.loadItems(locale)
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].URL = client.getLinkForEntry(items[i], locale)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	return items, nil
}
