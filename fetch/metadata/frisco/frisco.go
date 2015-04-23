package frisco

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
)

const IBUProperty = "IBU"
const ProfileURLProperty = "friscoProfileUrl"
const ServingSizeProperty = "friscoServingSize"
const FriscoDescription = "friscoDescription"

var ErrNoFriscoProfile = errors.New("no frisco profile")

func Agent() *httpagent.Agent {
	return httpagent.Win1252Agent()
}

func IsFrisco(bev model.Beverage) bool {
	return bev.Attribute(ProfileURLProperty) != ""
}

func FetchMetadata(bev model.Beverage) error {
	friscoProfile := bev.Attribute(ProfileURLProperty)
	if friscoProfile == "" {
		return ErrNoFriscoProfile
	}

	profileDoc, err := Agent().GetDoc(friscoProfile)
	if err != nil {
		return err
	}

	setFriscoMetadata(bev, profileDoc)
	return nil
}

var rBrewery = regexp.MustCompile(`Brewery: (.*)`)
var rName = regexp.MustCompile(`Name: (.*)`)
var rAbv = regexp.MustCompile(`ABV: (.*)%`)
var rServing = regexp.MustCompile(`Serving Size: (.*)`)
var rDesc = regexp.MustCompile(`(?s)Description: (.*)`)
var rType = regexp.MustCompile(`Description: (.*?)-`)
var rIBU = regexp.MustCompile(`(?s)Description: .*IBU: (\d+)`)
var rDescTypePrefix = regexp.MustCompile("^(.*?)-")
var rDescIBUPrefix = regexp.MustCompile(`^\s*IBU: \d+\s*-`)

func setFriscoMetadata(bev model.Beverage, doc *goquery.Document) {
	desc := doc.Find("[data-role='page'] [data-role='content']").Text()

	overwrite := bev.AccuracyScore() <= 2
	if bev.AccuracyScore() < 2 {
		bev.SetAccuracyScore(2)
	}
	bev.SetNeedSync(true)

	setExtractedText := func(reg *regexp.Regexp, current string, action func(string)) {
		if overwrite || current == "" {
			match := reg.FindStringSubmatch(desc)
			if match != nil && match[1] != "" {
				action(text.Normalize(match[1]))
			}
		}
	}

	setExtractedFloat64 := func(reg *regexp.Regexp, current float64, action func(float64)) {
		fakeVal := ""
		if current > 0.0 {
			fakeVal = "_"
		}
		setExtractedText(reg, fakeVal, func(abv string) {
			fabv, err := strconv.ParseFloat(abv, 64)
			if err == nil && fabv > 0.0 {
				action(fabv)
			}
		})
	}

	setExtractedText(rBrewery, bev.Brewer(), bev.SetBrewer)
	setExtractedText(rName, bev.Name(), bev.SetName)
	setExtractedText(rType, bev.Type(), bev.SetType)
	setExtractedFloat64(rAbv, bev.Abv(), bev.SetAbv)
	setExtractedText(rServing,
		bev.Attribute(ServingSizeProperty), func(serving string) {
			bev.SetAttribute(ServingSizeProperty, serving)
		})

	setExtractedText(rIBU,
		bev.Attribute(IBUProperty), func(ibu string) {
			bev.SetAttribute(IBUProperty, ibu)
		})

	descMatch := rDesc.FindStringSubmatch(desc)
	if descMatch != nil {
		exDesc := descMatch[1]
		exDesc = rDescTypePrefix.ReplaceAllLiteralString(exDesc, "")
		exDesc = rDescIBUPrefix.ReplaceAllLiteralString(exDesc, "")
		exDesc = text.NormalizeMultiline(exDesc)
		bev.SetAttribute(FriscoDescription, exDesc)

		if exDesc != "" && (overwrite || bev.Description() == "") {
			bev.SetDescription(exDesc)
		}
	}
}
