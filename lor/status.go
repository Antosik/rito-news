package lor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Antosik/rito-news/internal/serverstatus"
)

// Legends of Runeterra server status entry.
type StatusEntry serverstatus.Entry

// A client that allows to get Legends of Runeterra server status.
//
// Source - https://status.riotgames.com/lor
type StatusClient struct {
	// Available regions:
	// americas, apac, europe
	Region string
}

func (client StatusClient) loadItems(locale string) ([]serverstatus.Entry, error) {
	url := fmt.Sprintf(
		"https://bacon.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)

	items, err := serverstatus.GetItems(url, locale)
	if err != nil {
		return nil, fmt.Errorf("can't get items: %w", err)
	}

	return items, nil
}

func (client StatusClient) getLinkForEntry(entry serverstatus.Entry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/lor?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.UID,
	)
}

// Available locales:
// en-US, de-DE, es-ES, es-MX, fr-FR, it-IT, ja-JP, ko-KR,
// pl-PL, pt-BR, ru-RU, th-TH, tr-TR, vi-VN, zh-TW.
func (client StatusClient) GetItems(locale string) ([]StatusEntry, error) {
	items, err := client.loadItems(locale)
	if err != nil {
		return nil, err
	}

	results := make([]StatusEntry, len(items))
	for index, item := range items {
		results[index] = StatusEntry{
			UID:         item.UID,
			Author:      item.Author,
			Date:        item.Date,
			Description: item.Description,
			Title:       item.Title,
			URL:         client.getLinkForEntry(items[index], locale),
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	return results, nil
}
