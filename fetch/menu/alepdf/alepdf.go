package alepdf

import (
	"regexp"
	"strconv"

	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/pdftext"
)

const ServingProperty = "ale_houseServingSize"

var SectionHeadingFont = pdftext.Font{Name: "Duke-Fill", Size: 22}
var BevTitleFont = pdftext.Font{Name: "Stag-Semibold", Size: 11}
var BevTypeMetaFont = pdftext.Font{Name: "Stag-Semibold", Size: 9}
var BevBodyFont = pdftext.Font{Name: "Stag-Book", Size: 9}

type section struct {
	name     string
	modifier func(string) string
}

func (s section) modifyName(name string) string {
	if s.modifier != nil {
		return s.modifier(name)
	}
	return name
}

var sectionMap = map[string]section{
	"Guest Drafts": section{name: "drafts"},
	"OLIVER BREWING CO.": section{"oliver", func(name string) string {
		return "Oliver " + name
	}},
}

type alehouseReader struct {
	section string
	bev     model.Beverage
	*pdftext.Scanner
}

func (r *alehouseReader) knownSection(name string) bool {
	_, ok := sectionMap[name]
	return ok
}

func (r *alehouseReader) newBeverage(frag *pdftext.Fragment) model.Beverage {
	if frag.Font != BevTitleFont {
		return nil
	}
	return model.CreateBeverage(sectionMap[r.section].modifyName(frag.Text()))
}

func (r *alehouseReader) emitBev() model.Beverage {
	bev := r.bev
	r.bev = nil
	return bev
}

func (r *alehouseReader) NextBeverage() (model.Beverage, error) {
	for {
		textFrag := r.NextText()
		if textFrag == nil {
			return r.emitBev(), nil
		}

		if textFrag.Font == SectionHeadingFont {
			r.section = textFrag.Text()

			if lastBev := r.emitBev(); lastBev != nil {
				return lastBev, nil
			}
			continue
		}

		if !r.knownSection(r.section) {
			continue
		}

		if newBev := r.newBeverage(textFrag); newBev != nil {
			if r.bev != nil {
				result := r.bev
				r.bev = newBev
				return result, nil
			}
			r.bev = newBev
			continue
		}

		if r.bev == nil {
			continue
		}

		r.setBevType(textFrag)
		r.setBevDescriptionABV(textFrag)
		r.setPours(textFrag)
	}
}

func (r *alehouseReader) setBevType(frag *pdftext.Fragment) {
	if frag.Font != BevTypeMetaFont || r.bev.Type() != "" {
		return
	}
	r.bev.SetType(frag.Text())
}

var rABVSuffix = regexp.MustCompile(`\s+(\d+(?:[.]\d+)?)%`)

func (r *alehouseReader) setBevDescriptionABV(frag *pdftext.Fragment) {
	if frag.Font != BevBodyFont {
		return
	}
	desc := frag.Text()
	if abvMatch := rABVSuffix.FindStringSubmatch(desc); abvMatch != nil {
		desc = rABVSuffix.ReplaceAllString(desc, "")
		if abv, err := strconv.ParseFloat(abvMatch[1], 64); err == nil && abv > 0.0 {
			r.bev.SetAbv(abv)
		}
	}
	r.bev.SetDescription(desc)
}

var rPourPrice = regexp.MustCompile(`(\w+)(?: -)? (\$\d+(?:[.]\d+)?)`)

func (r *alehouseReader) setPours(frag *pdftext.Fragment) {
	if frag.Font != BevTypeMetaFont {
		return
	}
	if pour := rPourPrice.FindStringSubmatch(frag.Text()); pour != nil {
		servings := r.bev.Attribute(ServingProperty)
		if servings != "" {
			servings += ", "
		}
		servings += pour[1] + ": " + pour[2]
		r.bev.SetAttribute(ServingProperty, servings)
	}
}

func Parse(pdfname string) ([]model.Beverage, error) {
	scanner, err := pdftext.NewFileScanner(pdfname)
	if err != nil {
		return nil, err
	}

	aleReader := &alehouseReader{Scanner: scanner}

	result := []model.Beverage{}
	for {
		bev, err := aleReader.NextBeverage()
		if err != nil {
			return nil, err
		}
		if bev == nil {
			break
		}
		result = append(result, bev)
	}
	return result, nil
}
