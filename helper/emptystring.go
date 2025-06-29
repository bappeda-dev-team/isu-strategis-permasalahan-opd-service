package helper

import "strings"

func EmptyStringIfNull(s string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return ""
	}
	return s
}
