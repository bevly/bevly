package alepdf

import (
	"fmt"
	"os"
	"testing"

	"github.com/bevly/bevly/model"
)

var bevTests = []struct {
	name    string
	bevtype string
	abv     float64
	desc    string
	pours   string
}{
	{
		name:    "Jailbreak - Desserted",
		bevtype: "Chocolate Coconut Porter",
		desc:    "Easy sipping body with a surprising but pleasantly long, rich finish. Lots of coconut aroma but a pretty traditional palate.",
		abv:     6.9,
		pours:   "5oz: $2.75, 10oz: $4.95, 16oz: $6.95, 23oz: $9.25, Gr: $27.95",
	},
	{
		name:    "Oliver Draft Punk",
		bevtype: "American IPA",
		abv:     7,
	},
	{
		name:    "Jailbreak - Cafe Kavorka",
		bevtype: "Porter",
		abv:     5.5,
		desc:    "A lightly roasted porter fermented with sweet cherries and a late addition of tart, black cherries. The result is a semisweet & slightly tart perfectly dark porter with loads of complexity & depth of \"avor.",
	},
}

func TestParse(t *testing.T) {
	bevs, err := Parse("test/menu.pdf")
	if err != nil {
		t.Errorf("Error parsing PDF: %s", err)
		return
	}

	const expectedBevCount = 32
	if len(bevs) != expectedBevCount {
		t.Errorf("len(bevs) = %d, want %d", len(bevs), expectedBevCount)
	}

	for i, bev := range bevs {
		fmt.Fprintf(os.Stderr, "%d) %s (%s)\n", i+1, bev.DisplayName(),
			bev.Type())
	}

	findBev := func(name string) model.Beverage {
		for _, bev := range bevs {
			if bev.DisplayName() == name {
				return bev
			}
		}
		return nil
	}

	for _, test := range bevTests {
		bev := findBev(test.name)
		if bev == nil {
			t.Errorf("can't find bev: %#v", test.name)
			continue
		}
		if test.bevtype != "" && test.bevtype != bev.Type() {
			t.Errorf("bev.Type() == %#v, want %#v",
				bev.Type(), test.bevtype)
		}
		if test.desc != "" && test.desc != bev.Description() {
			t.Errorf("bev.Description() == %#v, want %#v",
				bev.Description(), test.desc)
		}
		if test.abv > 0 && test.abv != bev.Abv() {
			t.Errorf("bev.Abv() == %#v, want %#v", bev.Abv(), test.abv)
		}
		if test.pours != "" && bev.Attribute(ServingProperty) != test.pours {
			t.Errorf("bev.Attribute(%#v) == %#v, want %#v",
				ServingProperty, bev.Attribute(ServingProperty), test.pours)
		}
	}
}
