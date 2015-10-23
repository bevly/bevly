package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripNonAlpha(t *testing.T) {
	assert.Equal(t, "Push Ob session",
		StripNonAlpha("Push Ob\"session\""), "push")
	assert.Equal(t, "Baird Suruga Bay Imperial IPA RateBeer",
		StripNonAlpha("Baird Suruga Bay Imperial IPA - RateBeer"), "baird")
}

func TestMatchConfidence(t *testing.T) {
	assert.True(t,
		NameMatchConfidence(
			"Baird 合資会社 Suruga Bay Imperial IPA",
			"Baird Suruga Bay Imperial IPA - RateBeer") > 0.25,
		"Suruga Bay")

	assert.False(t,
		NameMatchConfidence(
			"Push Ob\"session\"",
			"Green Flash Hop Head Red Ale - Beer Advocate") > 0.25,
		"push")
	assert.True(t,
		NameMatchConfidence(
			"Lagunitas Sucks",
			"Lagunitas Sucks (Brown ...") > 0.25,
		"sucks")
	assert.True(t,
		NameMatchConfidence(
			"marstons pedigree CASK",
			"Pedigree | Marston, Thompson & Evershed, Plc. | Burton-on ...") > 0.25,
		"marston")
}

func TestStripNonNumeric(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"5.3%", "5.3"},
		{"a-53c", "-53"},
	}
	for _, test := range tests {
		actual := StripNonNumeric(test.in)
		if actual != test.expected {
			t.Errorf("StripNonNumeric(%#v) == %#v, want %#v", test.in, actual, test.expected)
		}
	}
}
