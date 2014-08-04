package beeradvocate

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/google"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoResults = errors.New("no results for beverage")
var ErrNotBABeer = errors.New("not a beer on BA")

func FetchMetadata(bev model.Beverage) error {
	bev.SetLink(google.SearchURL(bev.DisplayName()))

	log.Printf("Searching for BA profile for %s", bev)
	baUrl, err := FindProfile(bev)
	if err != nil {
		log.Printf("BA profile error for %s: %s", bev, err)
		return err
	}
	return fetchBAMetadata(bev, baUrl)
}

func FindProfile(bev model.Beverage) (string, error) {
	baUrl, err := baGoogle("site:beeradvocate.com " + bev.DisplayName())
	if err != nil {
		return "", err
	}
	if baUrl == "" {
		baUrl, err = baGoogle("beeradvocate " + bev.DisplayName())
		if err != nil {
			return "", err
		}
	}
	if baUrl == "" {
		return "", ErrNoResults
	}
	return baUrl, nil
}

func baGoogle(search string) (string, error) {
	results, err := google.Search(search)
	if err != nil {
		return "", err
	}
	if len(results) > 0 {
		for _, result := range results {
			urlString := result.URL.String()
			if strings.Contains(urlString, "beeradvocate.com/beer") {
				return urlString, nil
			}
		}
	}
	return "", nil
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
