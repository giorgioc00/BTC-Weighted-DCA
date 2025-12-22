package util

import (
	"errors"
	"strings"
	"time"
)

var dateLayouts = []string{
	"2006-01-02",
	"01/02/2006",
	"02/01/2006",
}

// ParseDateFlexible tries multiple layouts to parse a date string
func ParseDateFlexible(dateStr string) (time.Time, error) {
	clean := strings.TrimSpace(strings.Trim(dateStr, "\""))
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, clean); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unable to parse date: " + dateStr)
}
