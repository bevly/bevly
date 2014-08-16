package frisco

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"regexp"
	"strconv"
)

const ProfileURLProperty = "friscoProfileUrl"
const ServingSizeProperty = "friscoServingSize"
const FriscoDescription = "friscoDescription"

var ErrNoFriscoProfile = errors.New("no frisco profile")

func Agent() *httpagent.HttpAgent {
	agent := httpagent.Agent()
	agent.ForceEncoding = "latin1"
	return agent
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

func setFriscoMetadata(bev model.Beverage, doc *goquery.Document) {
	desc := doc.Find("[data-role='page'] [data-role='content']").Text()

	overwrite := bev.AccuracyScore() <= 2
	if bev.AccuracyScore() < 2 {
		bev.SetAccuracyScore(2)
	}
	bev.SetNeedSync(true)

	setExtractedText(rBrewery, desc, overwrite, bev.Brewer(), bev.SetBrewer)
	setExtractedText(rName, desc, overwrite, bev.Name(), bev.SetName)
	setExtractedFloat64(rAbv, desc, overwrite, bev.Abv(), bev.SetAbv)
	setExtractedText(rServing, desc, overwrite,
		bev.Attribute(ServingSizeProperty), func(serving string) {
			bev.SetAttribute(ServingSizeProperty, serving)
		})
	setExtractedText(rDesc, desc, true, bev.Description(),
		func(desc string) {
			if overwrite {
				bev.SetDescription(desc)
			}
			bev.SetAttribute(FriscoDescription, desc)
		})
}

func setExtractedText(reg *regexp.Regexp, haystack string,
	overwrite bool, current string, action func(string)) {
	if overwrite || current == "" {
		match := reg.FindStringSubmatch(haystack)
		if match != nil && match[1] != "" {
			action(text.Normalize(match[1]))
		}
	}
}

func setExtractedFloat64(reg *regexp.Regexp, text string, overwrite bool, current float64, action func(float64)) {
	fakeVal := ""
	if current > 0.0 {
		fakeVal = "_"
	}
	setExtractedText(reg, text, overwrite, fakeVal, func(abv string) {
		fabv, err := strconv.ParseFloat(abv, 64)
		if err == nil && fabv > 0.0 {
			action(fabv)
		}
	})
}
