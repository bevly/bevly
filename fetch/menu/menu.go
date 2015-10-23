package menu

import (
	"errors"
	"log"

	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/google"
)

var ErrEmptyMenu = errors.New("empty menu")

type menuFetcher func(model.MenuProvider) ([]model.Beverage, error)

var menuFetcherRegistry = map[string]menuFetcher{}

// Get list of beverages for a menu provider
func FetchMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	fetcher := menuFetcherRegistry[provider.MenuFormat()]
	if fetcher != nil {
		log.Printf("FetchMenu(%s): start fetch:%s\n",
			provider.ID(), provider.MenuFormat())
		beverages, err := fetcher(provider)
		if err != nil {
			return nil, err
		}

		search := google.DefaultSearch()
		for _, bev := range beverages {
			link := search.SearchURL(bev.DisplayName())
			bev.SetAttribute("googleLink", link)
			if bev.Link() == "" {
				bev.SetLink(link)
			}
		}
		return beverages, nil
	}

	log.Printf("FetchMenu(%s): no fetcher for %s", provider.ID(), provider.MenuFormat())
	return []model.Beverage{}, nil
}
