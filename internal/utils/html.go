package utils

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReadHTML(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("can't read html: %w", err)
	}

	return doc, nil
}
