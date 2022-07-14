package val

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Antosik/rito-news/internal/serverstatus"
)

// VALORANT server status entry
type StatusEntry serverstatus.Entry

// A client that allows to get VALORANT server status.
//
// Source - https://status.riotgames.com/valorant
type StatusClient struct {
	// Available regions:
	// ap, br, eu, kr, latam, na, pbe
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

// Available locales:
// en-US, ar-AE, de-DE, es-ES, es-MX, fr-FR, id-ID, it-IT, ja-JP,
// ko-KR, pl-PL, pt-BR, ru-RU, tr-TR, th-TH, vi-VN, zh-TW
func (client StatusClient) GetItems(locale string) ([]StatusEntry, error) {
	items, err := client.loadItems(locale)
	if err != nil {
		return nil, err
	}

	results := make([]StatusEntry, len(items))
	for i, item := range items {
		results[i] = StatusEntry{
			UID:         item.UID,
			Author:      item.Author,
			Date:        item.Date,
			Description: item.Description,
			Title:       item.Title,
			URL:         client.getLinkForEntry(items[i], locale),
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})

	return results, nil
}
