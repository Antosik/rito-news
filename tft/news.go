package tft

import (
	"github.com/Antosik/rito-news/internal/nextjsnews"
)

// Teamfight Tactics news entry
type NewsEntry = nextjsnews.Item

// A client that allows to get official Teamfight Tactics news.
//
// Source - https://teamfighttactics.leagueoflegends.com/
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, de-de, es-es, fr-fr, it-it, en-au, pl-pl, ru-ru,
	// el-gr, ro-ro, hu-hu, cs-cz, es-mx, pt-br, tr-tr, ko-kr, ja-jp
	// en-sg, en-ph, zh-tw, vi-vn, th-th
	Locale string
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	parser := nextjsnews.Parser{BaseURL: "https://teamfighttactics.leagueoflegends.com", Locale: client.Locale}

	items, err := parser.GetItems(count)
	if err != nil {
		return nil, err
	}

	return items, nil
}
