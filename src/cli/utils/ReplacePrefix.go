package utils

import (
	"strings"
)

func ReplacePrefix(s, prefix, newPrefix string) string {
	if strings.HasPrefix(s, prefix) {
		// TrimPrefix removes the prefix and then we add the newPrefix.
		return newPrefix + strings.TrimPrefix(s, prefix)
	}
	return s
}
