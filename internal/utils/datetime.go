package utils

import (
	"fmt"
	"time"
)

func ParseDateTimeWithLayouts(datetime string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		date, err := time.Parse(layout, datetime)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("can't parse datetime %s", datetime)
}
