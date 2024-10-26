package val

import (
	"github.com/Antosik/rito-news/internal/nextjsnews"
)

// VALORANT news entry
type NewsEntry = nextjsnews.Item

// A client that allows to get official VALORANT news.
//
// Source - https://playvalorant.com/en-us/news/
type NewsClient struct {
	// Available locales:
	// en-us, en-gb, de-de, es-es, fr-fr, it-it, pl-pl, ru-ru, tr-tr,
	// es-mx, id-id, ja-jp, ko-kr, pt-br, th-th, vi-vn, zh-tw, ar-ae
	Locale string
}

func (client NewsClient) GetItems(count int) ([]NewsEntry, error) {
	parser := nextjsnews.Parser{BaseURL: "https://playvalorant.com", Locale: client.Locale}

	items, err := parser.GetItems(count)
	if err != nil {
		return nil, err
	}

	return items, nil
}
