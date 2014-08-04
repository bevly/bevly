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

	profileDoc, err := httpagent.Agent().GetDoc(friscoProfile)
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
	setExtractedText(rBrewery, desc, bev.SetBrewer)
	setExtractedText(rName, desc, bev.SetName)
	setExtractedFloat64(rAbv, desc, bev.SetAbv)
	setExtractedText(rServing, desc, func(serving string) {
		bev.SetAttribute(ServingSizeProperty, serving)
	})
	setExtractedText(rDesc, desc, bev.SetDescription)
}

func setExtractedText(reg *regexp.Regexp, haystack string, action func(string)) {
	match := reg.FindStringSubmatch(haystack)
	if match != nil {
		action(text.Normalize(match[1]))
	}
}

func setExtractedFloat64(reg *regexp.Regexp, text string, action func(float64)) {
	setExtractedText(reg, text, func(abv string) {
		fabv, err := strconv.ParseFloat(abv, 64)
		if err == nil && fabv > 0.0 {
			action(fabv)
		}
	})
}
