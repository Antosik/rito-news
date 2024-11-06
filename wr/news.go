package wr

import (
	"fmt"

	"github.com/Antosik/rito-news/internal/nextjsnews"
)

// Wild Rift news entry.
type NewsEntry = nextjsnews.Item

// A client that allows to get official Wild Rift news.
//
// Source - https://wildrift.leagueoflegends.com/en-us/news/
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, fr-fr, de-de, es-es, it-it, pl-pl, ru-ru, tr-tr, id-id,
	// ms-my, pt-br, ja-jp, ko-kr, zh-tw, th-th, vi-vn, es-mx, en-sg, ar-ae
	Locale string
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	parser := nextjsnews.Parser{BaseURL: "https://wildrift.leagueoflegends.com", Locale: client.Locale}

	items, err := parser.GetItems(count)
	if err != nil {
		return nil, fmt.Errorf("can't get news items: %w", err)
	}

	return items, nil
}
