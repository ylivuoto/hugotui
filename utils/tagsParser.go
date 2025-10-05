package utils

import "strings"

func ParseTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	formatted := make([]string, len(tags))
	for i, t := range tags {
		t = strings.TrimSpace(t)
		if !strings.HasPrefix(t, "#") {
			t = "#" + t
		}
		formatted[i] = t
	}
	return strings.Join(formatted, " ")
}
