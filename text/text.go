package text

import (
	"regexp"
	"strings"
)

var rMultispaceRegex = regexp.MustCompile(`\s{2,}`)
var rWhitespace = regexp.MustCompile(`\s`)

func Normalize(name string) string {
	return rMultispaceRegex.ReplaceAllString(strings.TrimSpace(name), " ")
}

var rNonLetter = regexp.MustCompile(`\PL+`)

func StripNonAlpha(text string) string {
	return Normalize(rNonLetter.ReplaceAllString(text, " "))
}

// Returns a confidence level that two names match. Given two names,
// splits them into a set of words, then returns the fraction of the
// words in A that are in the set multiplied by the fraction of the
// words in B that are in the set.
//
// Empty strings never match any other string, including other empty
// strings.
func NameMatchConfidence(a, b string) float64 {
	wordsA := SplitWords(strings.ToLower(StripNonAlpha(a)))
	wordsB := SplitWords(strings.ToLower(StripNonAlpha(b)))

	if len(wordsA) == 0 || len(wordsB) == 0 {
		return 0
	}

	nIntersect := 0
	for _, wordA := range wordsA {
		for _, wordB := range wordsB {
			if wordA == wordB {
				nIntersect++
				break
			}
		}
	}
	return float64(nIntersect*nIntersect) / float64(len(wordsA)*len(wordsB))
}

func SplitWords(text string) []string {
	return rWhitespace.Split(text, -1)
}
