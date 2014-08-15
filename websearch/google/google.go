package google

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/text"
	"github.com/bevly/bevly/throttle"
	"github.com/bevly/bevly/websearch"
	"log"
	"net/url"
)

const SearchBaseURL = "https://www.google.com/search"

var GoogleThrottle = throttle.Default(SearchBaseURL)

type GoogleSearch struct {
	BaseURL  string
	Throttle *throttle.Throttle
}

func DefaultSearch() websearch.Search {
	return &GoogleSearch{
		BaseURL:  SearchBaseURL,
		Throttle: GoogleThrottle,
	}
}

func SearchWithURL(baseURL string) websearch.Search {
	return &GoogleSearch{BaseURL: baseURL}
}

func (g *GoogleSearch) Search(terms string) ([]websearch.Result, error) {
	if g.Throttle != nil {
		g.Throttle.DelayInvocation()
	}

	url := g.SearchURL(terms)
	log.Printf("GoogleSearch(%s): GET %s", terms, url)
	response, err := httpagent.Agent().Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	results := []websearch.Result{}
	results = addBodySearchResults(doc, results)

	log.Printf("GoogleSearch(%s): %d results\n", terms, len(results))
	for i, result := range results {
		log.Printf("%2d) %s (%s)\n", i+1, result.Text, result.URL)
	}
	return results, nil
}

func addBodySearchResults(doc *goquery.Document, results []websearch.Result) []websearch.Result {
	doc.Find("h3.r").Each(func(i int, s *goquery.Selection) {
		anchor := s.Find("a").First()
		href, exists := anchor.Attr("href")
		if exists {
			results = append(results, websearch.Result{
				URL:  href,
				Text: text.Normalize(anchor.Text()),
			})
		}
	})
	return results
}

func (g *GoogleSearch) SearchURL(terms string) string {
	return g.BaseURL + "?" + url.Values{"q": {terms}}.Encode()
}
