package beeradvocate

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"github.com/bevly/bevly/websearch"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoResults = errors.New("no results for beverage")
var ErrNotBABeer = errors.New("not a beer on BA")

func FetchMetadata(bev model.Beverage, search websearch.Search) error {
	bev.SetLink(search.SearchURL(bev.DisplayName()))

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
		return "", ErrNoResults
	}
	return baUrl, nil
}

func baSearch(name string, search websearch.Search) (string, error) {
	terms := "beeradvocate " + name
	results, err := search.Search(terms)
	if err != nil {
		return "", err
	}
	log.Printf("baSearch(%s): %d results\n", search, len(results))
	if len(results) > 0 {
		for _, result := range results {
			urlString := result.URL
			log.Printf("baSearch(%s): considering %s (%s)\n",
				search, result.Text, result.URL)
			if IsBeerAdvocateProfile(urlString) {
				confidence := text.NameMatchConfidence(result.Text, name)
				if confidence < 0.2 {
					log.Printf("baSearch(%s): rejecting %s (confidence: %.2f%%)\n",
						search, result.Text, confidence)
				}
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

	bev.SetLink(metaURL)
	setBATitleBrewer(bev, doc)
	setBATypeAbv(bev, doc)
	setBARatings(bev, doc)
	return nil
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
