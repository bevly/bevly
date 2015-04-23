package pdftext

import (
	"fmt"
	"os"
	"testing"
)

func TestDump(t *testing.T) {
	scanner, err := NewFileScanner("test/menu.pdf")
	if err != nil {
		t.Errorf("Failed to create scanner: %s", err)
		return
	}
	for {
		text := scanner.NextText()
		if text == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Fragment: %#v (%s)\n", text, text.Text())
	}
}

func TestNextText(t *testing.T) {
	scanner, err := NewFileScanner("test/menu.pdf")
	if err != nil {
		t.Errorf("Failed to create scanner: %s", err)
		return
	}

	expected := []string{
		"Guest Drafts",
		"Jailbreak - Desserted",
		"Chocolate Coconut Porter",
		"Easy sipping body with a surprising but pleasantly long, rich finish. Lots of coconut aroma but a pretty traditional palate. 6.9%",
		"5oz - $2.75",
		"|",
		"10oz - $4.95",
		"|",
		"16oz - $6.95",
		"|",
		"23oz - $9.25",
		"|",
		"Gr $27.95",

		"Jailbreak - Big Punisher",
		"Double IPA",
		"A well-balanced double IPA with a semisweet malt backbone and complimented with generous amounts of citrus & tropical fruit hops. Rich, delicious, and rewardingly punishing. 8.5%",
		"5oz - $2.75",
		"|",
		"10oz - $4.95",
		"|",
		"16oz - $6.95",
		"|",
		"23oz - $9.25",
		"|",
		"Gr $27.95",
	}

	mustSee := []string{
		"Evil Twin - I Love You With My Stout",
		"Draft Punk",
	}

	scannedFragments := []*Fragment{}
	seenText := map[string]struct{}{}
	for {
		next := scanner.NextText()
		if next == nil {
			break
		}
		scannedFragments = append(scannedFragments, next)
		seenText[next.Text()] = struct{}{}
	}

	for i, expectedFragment := range expected {
		if i >= len(scannedFragments) {
			t.Errorf("NextText()#%d == nil, want %#v", i, expectedFragment)
			continue
		}
		actual := scannedFragments[i]
		if actual.Text() != expectedFragment {
			t.Errorf("NextText()#%d == %#v, want %#v", i, actual.Text(), expectedFragment)
		}
	}

	for _, see := range mustSee {
		if _, ok := seenText[see]; !ok {
			t.Errorf("expected to see %#v, but did not find it", see)
		}
	}
}
