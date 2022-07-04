package serverstatus

import (
	"strconv"
	"strings"
)

func getLocaleFromServerStatusEntryTranslations(translations []serverStatusEntryAPITranslation, locale string) string {
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

func transformServerStatusEntryToNewsItems(status serverStatusAPIEntry, locale string) []ServerStatusEntry {
	items := make([]ServerStatusEntry, 0, len(status.Updates))

	title := getLocaleFromServerStatusEntryTranslations(status.Titles, locale)

	for _, update := range status.Updates {
		if !update.Publish {
			continue
		}

		items = append(items, ServerStatusEntry{
			UID:         strconv.Itoa(update.Id),
			Title:       title,
			Description: getLocaleFromServerStatusEntryTranslations(update.Translations, locale),
			Author:      update.Author,
			Date:        update.CreatedAt,
		})
	}

	return items
}

func TransformServerStatusToNewsItems(status serverStatusAPIResponse, locale string) []ServerStatusEntry {
	statuses := make([]serverStatusAPIEntry, 0, len(status.Incidents)+len(status.Maintenances))
	statuses = append(statuses, status.Incidents...)
	statuses = append(statuses, status.Maintenances...)

	var items []ServerStatusEntry

	for _, entry := range statuses {
		items = append(items, transformServerStatusEntryToNewsItems(entry, locale)...)
	}

	return items
}
