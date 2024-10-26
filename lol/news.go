package lol

import (
	"github.com/Antosik/rito-news/internal/nextjsnews"
)

// League of Legends news entry
type NewsEntry = nextjsnews.Item

// A client that allows to get official League of Legends news.
//
// Source - https://www.leagueoflegends.com/en-us/news
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, de-de, es-es, fr-fr, it-it, en-pl, pl-pl, el-gr, ro-ro,
	// hu-hu, cs-cz, es-mx, pt-br, ja-jp, ru-ru, tr-tr, en-au, ko-kr,
	// en-sg, en-ph, vi-vn, th-th, zh-tw
	Locale string
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	parser := nextjsnews.Parser{BaseURL: "https://www.leagueoflegends.com", Locale: client.Locale}

	items, err := parser.GetItems(count)
	if err != nil {
		return nil, err
	}

	return items, nil
}
