package utils

import (
	"errors"
	"fmt"
	"time"
)

var errDatetimeParse = errors.New("can't parse datetime")

func ParseDateTimeWithLayouts(datetime string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		date, err := time.Parse(layout, datetime)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("%w: no valid layouts found for %s", errDatetimeParse, datetime)
}
