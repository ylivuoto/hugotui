package utils

import (
	"time"
)

func FormatHugoDate(dateStr string) string {
	// try full RFC3339 first (Hugo default)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// fall back to plain date only (some posts use YYYY-MM-DD)
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			// return unchanged if nothing works
			return dateStr
		}
	}
	// now format nicely
	return t.Format("01.02.2006 15:04")
}
