package menu

import (
	"errors"
	"github.com/bevly/bevly/model"
	"log"
)

var ErrEmptyMenu = errors.New("empty menu")

type menuFetcher func(model.MenuProvider) ([]model.Beverage, error)

var menuFetcherRegistry = map[string]menuFetcher{}

// Get list of beverages for a menu provider
func FetchMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	fetcher := menuFetcherRegistry[provider.MenuFormat()]
	if fetcher != nil {
		log.Printf("FetchMenu(%s): start fetch:%s\n",
			provider.Id(), provider.MenuFormat())
		return fetcher(provider)
	}

	log.Printf("FetchMenu(%s): no fetcher for %s", provider.Id(), provider.MenuFormat())
	return []model.Beverage{}, nil
}
