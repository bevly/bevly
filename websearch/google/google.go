package google

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/text"
	"github.com/bevly/bevly/websearch"
	"net/url"
)

const SearchBaseURL = "https://www.google.com/search"

type GoogleSearch struct {
	BaseURL string
}

func DefaultSearch() websearch.Search {
	return &GoogleSearch{BaseURL: SearchBaseURL}
}

func SearchWithURL(baseURL string) websearch.Search {
	return &GoogleSearch{BaseURL: baseURL}
}

func (g *GoogleSearch) Search(terms string) ([]websearch.Result, error) {
	response, err := httpagent.Agent().Get(g.SearchURL(terms))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	results := []websearch.Result{}
	doc.Find(".srg li.g").Each(func(i int, s *goquery.Selection) {
		anchor := s.Find(".r a")
		href, exists := anchor.Attr("href")
		if exists {
			if err == nil {
				results = append(results, websearch.Result{
					URL:  href,
					Text: text.Normalize(anchor.Text()),
				})
			}
		}
	})
	return results, nil
}

func (g *GoogleSearch) SearchURL(terms string) string {
	return g.BaseURL + "?" + url.Values{"q": {terms}}.Encode()
}
