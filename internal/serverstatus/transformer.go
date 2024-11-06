package serverstatus

import (
	"strconv"
	"strings"
)

func getLocaleFromRawTranslationEntry(translations []rawTranslationEntry, locale string) string {
	var (
		fallback string
		result   string
	)

	for _, translation := range translations {
		if strings.EqualFold(translation.Locale, locale) ||
			strings.EqualFold(translation.Locale, strings.ReplaceAll(locale, "-", "_")) {
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

func transformRawEntryToEntry(status rawEntry, locale string) []Entry {
	items := make([]Entry, 0, len(status.Updates))

	title := getLocaleFromRawTranslationEntry(status.Titles, locale)

	for _, update := range status.Updates {
		if !update.Publish {
			continue
		}

		items = append(items, Entry{
			UID:         strconv.Itoa(update.ID),
			Title:       title,
			Description: getLocaleFromRawTranslationEntry(update.Translations, locale),
			Author:      update.Author,
			Date:        update.CreatedAt,
			URL:         "",
		})
	}

	return items
}

func transformRawResponseToEntry(status rawResponse, locale string) []Entry {
	statuses := make([]rawEntry, 0, len(status.Incidents)+len(status.Maintenances))
	statuses = append(statuses, status.Incidents...)
	statuses = append(statuses, status.Maintenances...)

	var items []Entry

	for _, entry := range statuses {
		items = append(items, transformRawEntryToEntry(entry, locale)...)
	}

	return items
}
