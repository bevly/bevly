package text

import (
	"regexp"
	"strings"
)

var whitespaceRegex = regexp.MustCompile(`\s{2,}`)

func Normalize(name string) string {
	return whitespaceRegex.ReplaceAllString(strings.TrimSpace(name), " ")
}
