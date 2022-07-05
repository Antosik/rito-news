package wr

import (
	"fmt"
	"rito-news/lib/serverstatus"
	"strings"
)

type StatusClient struct {
	Region string
}

func (client StatusClient) loadItems(locale string) ([]serverstatus.Entry, error) {
	url := fmt.Sprintf(
		"https://wildrift.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)

	return serverstatus.GetItems(url, locale)
}

func (client StatusClient) getLinkForEntry(entry serverstatus.Entry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/wildrift?region=%s&locale=%s&id=%s",
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

	return items, nil
}
