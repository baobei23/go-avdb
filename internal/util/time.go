package util

import (
	"strings"
	"time"
)

func ParseDate(dateStr string) *time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return nil
	}

	layouts := []string{
		"2006-01-02",
		"2 Jan, 2006",
		"02 Jan, 2006",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return &t
		}
	}

	return nil // Default fallback
}
