package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReadHTML(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
