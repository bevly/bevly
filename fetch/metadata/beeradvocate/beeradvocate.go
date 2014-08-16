package beeradvocate

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	bevHtml "github.com/bevly/bevly/html"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"github.com/bevly/bevly/websearch"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoResults = errors.New("no results for beverage")
var ErrNotBABeer = errors.New("not a beer on BA")

const DescriptionProperty = "baDescription"

func FetchMetadata(bev model.Beverage, search websearch.Search) error {
	log.Printf("Searching for BA profile for %s", bev)
	baUrl, err := FindProfile(bev, search)
	if err != nil {
		log.Printf("BA profile error for %s: %s", bev, err)
		return err
	}
	return fetchBAMetadata(bev, baUrl)
}

func FindProfile(bev model.Beverage, s websearch.Search) (string, error) {
	baUrl, err := baSearch(bev.DisplayName(), s)
	if err != nil {
		return "", err
	}
	if baUrl == "" {
		// No results, but we don't want to resync repeatedly anyway:
		bev.SetNeedSync(true)
		return "", ErrNoResults
	}
	return baUrl, nil
}

func stripBAName(text string) string {
	return strings.Replace(text, "Beer Advocate", "", 1)
}

func baSearch(name string, search websearch.Search) (string, error) {
	terms := "beeradvocate " + name
	results, err := search.Search(terms)
	if err != nil {
		return "", err
	}
	log.Printf("baSearch(%s): %d results\n", terms, len(results))
	if len(results) > 0 {
		for _, result := range results {
			urlString := result.URL
			log.Printf("baSearch(%s): considering %s (%s)\n",
				terms, result.Text, result.URL)
			if IsBeerAdvocateProfile(urlString) {
				cleansedText := stripBAName(result.Text)
				confidence := text.NameMatchConfidence(cleansedText, name)
				if confidence < 0.13 {
					log.Printf("baSearch(%s): rejecting %s (confidence: %.2f%%)\n",
						terms, cleansedText, confidence*100)
					continue
				}
				log.Printf("baSearch(%s): accepting %s (confidence: %.2f%%)\n",
					terms, cleansedText, confidence*100)
				return urlString, nil
			}
		}
	}
	return "", nil
}

var rBeerAdvocateProfileURL = regexp.MustCompile(`beeradvocate.*?/beer/profile/\d+/\d+`)

func IsBeerAdvocateProfile(url string) bool {
	return rBeerAdvocateProfileURL.FindString(url) != ""
}

func fetchBAMetadata(bev model.Beverage, metaURL string) error {
	log.Printf("fetchBAMetadata(%s, %s)", bev, metaURL)
	doc, err := httpagent.Agent().GetDoc(metaURL)
	if err != nil {
		log.Printf("fetchBAMetadata(%s, %s) failed: %s", bev, metaURL, err)
		return err
	}
	if !isBABeer(doc) {
		log.Printf("fetchBAMetadata: Ignoring %s: does not seem to be beer\n",
			metaURL)
		return ErrNotBABeer
	}

	bev.SetNeedSync(true)
	bev.SetAccuracyScore(10)
	bev.SetLink(metaURL)
	setBATitleBrewer(bev, doc)
	setBATypeAbv(bev, doc)
	setBARatings(bev, doc)
	setBADescription(bev, doc)
	return nil
}

var rDescRegexp = regexp.MustCompile(`(?s)Notes/Commercial Description:</b>(.*)`)

func setBADescription(bev model.Beverage, doc *goquery.Document) {
	parentHtml, err := doc.Find("a[href^=\"/beer/style\"]").First().Parent().Html()
	if err != nil {
		log.Printf("setBADescription(%s): error getting description: %s\n",
			bev, err)
		return
	}

	match := rDescRegexp.FindStringSubmatch(parentHtml)
	if match == nil {
		log.Printf("setBADescription(%s): no description heading?\n", bev)
		return
	}

	desc := cleanseBADescription(match[1])
	if desc != "" && desc != "No notes at this time." {
		log.Printf("setBADescription(%s): desc=%s\n", bev, desc)
		bev.SetDescription(desc)
		bev.SetAttribute(DescriptionProperty, desc)
	} else {
		log.Printf("setBADescription(%s): no desc\n", bev)
	}
}

var rAddedBy = regexp.MustCompile(`\(Beer added by:.*`)

func cleanseBADescription(desc string) string {
	desc = rAddedBy.ReplaceAllLiteralString(
		strings.Replace(desc, "\n", " ", -1), "")
	return bevHtml.Text(html.UnescapeString(
		text.NormalizeMultiline(strings.Replace(desc, "<br/>", "\n", -1))))
}

func setBARatings(bev model.Beverage, doc *goquery.Document) {
	scores := doc.Find(".BAscore_big")
	if scores.Size() < 1 {
		log.Printf("setBARatings: can't find ratings for %s\n", bev)
		return
	}
	addRating(bev, "BA", scores.First().Text())
	addRating(bev, "BAbro", scores.Last().Text())
}

func addRating(bev model.Beverage, name string, text string) {
	ratingPerc, err := strconv.ParseInt(text, 10, 32)
	if err == nil && ratingPerc > 0 {
		bev.AddRating(model.CreateRating(name, int(ratingPerc)))
	}
}

func setBATitleBrewer(bev model.Beverage, doc *goquery.Document) {
	header := doc.Find(".titleBar h1")
	combinedTitle := header.Text()
	breweryTitle := header.Find("span").Text()
	beerName := strings.Replace(combinedTitle, breweryTitle, "", 1)

	bev.SetName(text.Normalize(beerName))
	bev.SetBrewer(normalizeBrewer(breweryTitle))
	log.Printf("setBATitleBrewer: %s name=%s brewer=%s\n", bev, bev.Name(), bev.Brewer())
}

func setBATypeAbv(bev model.Beverage, doc *goquery.Document) {
	styleTag := doc.Find("a[href^=\"/beer/style\"]").First()
	style := styleTag.Text()
	bev.SetType(style)

	rawInfoText := styleTag.Parent().Text()
	abv := extractAbv(rawInfoText)
	bev.SetAbv(abv)
}

var abvPattern = regexp.MustCompile(`(?s)ABV.*?([\d.]+)%`)

func extractAbv(text string) float64 {
	match := abvPattern.FindStringSubmatch(text)
	if match != nil {
		abv, err := strconv.ParseFloat(match[1], 64)
		if err == nil {
			return abv
		}
		log.Printf("extractAbv: failed to parse ABV:%s\n", match[1])
	}
	log.Printf("extractAbv: no match for ABV pattern\n")
	return 0.0
}

var brewerLeadingDash = regexp.MustCompile(`^\s*-\s*`)

func normalizeBrewer(brewer string) string {
	return text.Normalize(brewerLeadingDash.ReplaceAllString(brewer, ""))
}

func isBABeer(doc *goquery.Document) bool {
	return doc.Find(".crust.selectedTabCrumb a span").First().Text() == "Beers"
}
