package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
}