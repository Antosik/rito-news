package lol

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Antosik/rito-news/internal/serverstatus"
)

// League of Legends server status entry
type StatusEntry serverstatus.Entry

// A client that allows to get League of Legends server status.
//
// Source - https://status.riotgames.com/lol;
type StatusClient struct {
	// Available regions and locales:
	// br1 (en-US, pt-BR); eun1 (en-US, en-GB, cs-CZ, el-GR, hu-HU, pl-PL, ro-RO); euw1 (en-US, en-GB, de-DE, es-ES, fr-FR, it-IT);
	// jp1 (en-US, ja-JP); kr1 (en-US, ko-KR); la1 (en-US, es-MX); la2 (en-US, es-AR);
	// na1 (en-US); oc1 (en-US, en-AU); ru1 (en-US, ru-RU); tr1 (en-US, tr-TR);
	// pbe (en-US, cs-CZ, de-DE, el-GR, es-MX, es-ES, fr-FR, hu-HU, it-IT, ja-JP, ko-KR, pl-PL, pt-BR, ro-RO, ru-RU, tr-TR).
	Region string
}

func (client StatusClient) loadItems(locale string) ([]serverstatus.Entry, error) {
	url := fmt.Sprintf(
		"https://lol.secure.dyn.riotcdn.net/channels/public/x/status/%s.json",
		client.Region,
	)

	return serverstatus.GetItems(url, locale)
}

func (client StatusClient) getLinkForEntry(entry serverstatus.Entry, locale string) string {
	return fmt.Sprintf(
		"https://status.riotgames.com/lol?region=%s&locale=%s&id=%s",
		client.Region,
		strings.ReplaceAll(locale, "-", "_"),
		entry.UID,
	)
}

//
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
