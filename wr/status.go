package wr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Antosik/rito-news/internal/serverstatus"
)

// Wild Rift server status entry
type StatusEntry serverstatus.Entry

// A client that allows to get Wild Rift server status.
//
// Source - https://status.riotgames.com/wildrift
type StatusClient struct {
	// Available regions:
	// br, eu, jp, kr, latam, mei, na, ru, sea
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

// Available locales:
// en-US, de-DE, en-GB, es-MX, es-ES, fr-FR, id-ID, it-IT, ja-JP, ko-KR,
// pl-PL, ms-MY, pt-BR, ru-RU, th-TH, tr-TR, vi-VN, zh-MY, zh-TW
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
