package text

import (
	"regexp"
	"strings"
)

var whitespaceRegex = regexp.MustCompile(`\s{2,}`)

func NormalizeName(name string) string {
	return whitespaceRegex.ReplaceAllString(strings.TrimSpace(name), " ")
}
