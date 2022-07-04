package serverstatus

import (
	"rito-news/utils/abstract"
	"strconv"
	"strings"
)

func getLocaleFromServerStatusEntryTranslations(translations []ServerStatusEntryTranslation, locale string) string {
	var fallback string
	var result string

	for _, translation := range translations {
		if strings.EqualFold(translation.Locale, locale) || strings.EqualFold(translation.Locale, strings.ReplaceAll(locale, "-", "_")) {
			result = translation.Content
			break
		}

		if strings.EqualFold(translation.Locale, "en_US") {
			fallback = translation.Content
		}
	}

	if result != "" {
		return result
	}

	return fallback
}

func transformServerStatusEntryToNewsItems(status ServerStatusEntry, locale string) []abstract.NewsItem {
	items := make([]abstract.NewsItem, len(status.Updates))

	title := getLocaleFromServerStatusEntryTranslations(status.Titles, locale)

	for i, update := range status.Updates {
		items[i] = abstract.NewsItem{
			Id:        strconv.Itoa(update.Id),
			Title:     title,
			Summary:   getLocaleFromServerStatusEntryTranslations(update.Translations, locale),
			Author:    update.Author,
			CreatedAt: update.CreatedAt,
			UpdatedAt: update.CreatedAt,
		}
	}

	return items
}

func TransformServerStatusToNewsItems(status ServerStatusResponse, locale string) []abstract.NewsItem {
	statuses := make([]ServerStatusEntry, 0, len(status.Incidents)+len(status.Maintenances))
	statuses = append(statuses, status.Incidents...)
	statuses = append(statuses, status.Maintenances...)

	var items []abstract.NewsItem

	for _, entry := range statuses {
		items = append(items, transformServerStatusEntryToNewsItems(entry, locale)...)
	}

	return items
}
