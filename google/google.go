package google

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/text"
	"net/url"
)

type Result struct {
	URL  *url.URL
	Text string
}

var SearchBaseURL = "https://www.google.com/search"

func Search(terms string) ([]Result, error) {
	response, err := httpagent.Agent().Get(SearchURL(terms))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	results := []Result{}
	doc.Find(".srg li.g").Each(func(i int, s *goquery.Selection) {
		anchor := s.Find(".r a")
		href, exists := anchor.Attr("href")
		if exists {
			parsedURL, err := url.Parse(href)
			if err == nil {
				results = append(results, Result{
					URL:  parsedURL,
					Text: text.NormalizeName(anchor.Text()),
				})
			}
		}
	})
	return results, nil
}

func SearchURL(terms string) string {
	return SearchBaseURL + "?" + url.Values{"q": {terms}}.Encode()
}
