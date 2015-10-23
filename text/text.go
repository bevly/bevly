package text

import (
	"math"
	"regexp"
	"strings"
)

var rMultispaceRegex = regexp.MustCompile(`\s+`)
var rWhitespace = regexp.MustCompile(`\s`)
var rPlainSpace = regexp.MustCompile(` +`)

func Normalize(name string) string {
	return rMultispaceRegex.ReplaceAllString(strings.TrimSpace(name), " ")
}

func NormalizeMultiline(name string) string {
	return rPlainSpace.ReplaceAllString(strings.TrimSpace(name), " ")
}

var rNonLetter = regexp.MustCompile(`\PL+`)

func StripNonAlpha(text string) string {
	return Normalize(rNonLetter.ReplaceAllString(text, " "))
}

var rNonDigitDot = regexp.MustCompile(`[^0-9.+-]`)

// StripNonNumeric removes all characters from text that are not in the set
// "0-9.+-"
func StripNonNumeric(text string) string {
	return rNonDigitDot.ReplaceAllString(text, "")
}

// NameMatchConfidence returns a confidence level that two names match. Given
// two names, splits them into a set of words, then returns the fraction of the
// words in A that are in the set multiplied by the fraction of the words in B
// that are in the set.
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
	return math.Max(
		float64(nIntersect)/float64(len(wordsA)),
		float64(nIntersect)/float64(len(wordsB)))
}

func SplitWords(text string) []string {
	return rWhitespace.Split(text, -1)
}
