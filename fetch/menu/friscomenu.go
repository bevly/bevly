package menu

import (
	"log"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
)

func init() {
	menuFetcherRegistry["frisco"] = friscoMenu
}

func friscoMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	agent := frisco.Agent()
	response, err := agent.Get(provider.URL())
	if err != nil {
		log.Printf("friscoMenu: Get(%s) failed: %s\n", provider.URL(), err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Printf("friscoMenu: Failed to create document from %s: %s\n", provider.URL(), err)
		return nil, err
	}
	beverages, err := friscoDrafts(doc, response.Request.URL)
	log.Printf("Frisco: parsed %d beverages from %s\n", len(beverages), provider.URL())
	return beverages, err
}

func friscoPourClassSize(pourClass string) (oz string) {
	switch pourClass {
	case "eightounce":
		return "8oz"
	case "tenounce":
		return "10oz"
	case "twelveounce":
		return "12oz"
	case "sixteenounce":
		return "16oz"
	default:
		return ""
	}
}

func friscoPourSize(s *goquery.Selection) (oz string) {
	pourNode := s.Children().First()
	if pourClass, exists := pourNode.Attr("class"); exists {
		return friscoPourClassSize(pourClass)
	}
	return ""
}

func parseABV(abv string) float64 {
	abv = text.StripNonNumeric(abv)
	fabv, err := strconv.ParseFloat(abv, 64)
	if err != nil {
		return 0
	}
	return fabv
}

func friscoDrafts(doc *goquery.Document, url *url.URL) ([]model.Beverage, error) {
	beers := []model.Beverage{}
	doc.Find(".row").Each(func(i int, s *goquery.Selection) {
		beerName := text.Normalize(s.Find(".name").Text())
		abv := s.Find(".abv").Text()
		pour := friscoPourSize(s)
		beer := model.CreateBeverage(beerName)
		beer.SetAttribute(frisco.ProfileURLProperty, url.String())
		if pour != "" {
			beer.SetAttribute(frisco.ServingSizeProperty, pour)
		}
		if abv != "" {
			beer.SetAbv(parseABV(abv))
		}
		beers = append(beers, beer)
	})
	if len(beers) == 0 {
		return nil, ErrEmptyMenu
	}
	return beers, nil
}
