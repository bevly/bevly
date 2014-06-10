package repository

import (
	"github.com/bevly/bevly/model"
	"log"
)

type Repository interface {
	MenuProviders() []model.MenuProvider
	ProviderById(id string) model.MenuProvider
	ProviderBeverages(provider model.MenuProvider) []model.Beverage
	ProviderIdBeverages(providerName string) []model.Beverage
	BeveragesNeedingSync() []model.Beverage

	// TODO
	SetBeverageMenu(provider model.MenuProvider, menu []model.Beverage)
	SaveBeverage(beverage model.Beverage)
	// BeverageByName(name string)

	// ClearMenus()
}

func DefaultRepository() Repository {
	return &stubRepository{}
}

type stubRepository struct{}

func (s *stubRepository) MenuProviders() []model.MenuProvider {
	return []model.MenuProvider{
		model.CreateMenuProvider("frisco", "Frisco", "http://beer.friscogrille.com/", "frisco"),
	}
}

func (s *stubRepository) ProviderById(id string) model.MenuProvider {
	log.Printf("Looking for provider named \"%s\"\n", id)
	for _, prov := range s.MenuProviders() {
		if prov.Id() == id {
			return prov
		}
	}
	return nil
}

func (s *stubRepository) ProviderBeverages(prov model.MenuProvider) []model.Beverage {
	if prov == nil {
		return []model.Beverage{}
	}
	return []model.Beverage{
		model.CreateBeverageAbvTypeRatingLink(
			"Mikkeller Beer Geek Brunch Weasel",
			10.9,
			"", 0, "",
			"http://www.google.com/search?q=Mikkeller+Beer+Geek+Brunch+Weasel"),

		model.CreateBeverageAbvTypeRatingLink(
			"Dogfish Head Olde School Barleywine",
			15.0,
			"Barleywine", 0, "",
			"http://www.google.com/search?q=Dogfish+Head+Olde+School+Barleywine"),

		model.CreateBeverageAbvTypeRatingLink(
			"Blue Point Toasted Lager",
			5.3,
			"", 80, "BA", "http://beeradvocate.com/beer/profile/764/2318/"),
	}
}

func (s *stubRepository) ProviderIdBeverages(provId string) []model.Beverage {
	return s.ProviderBeverages(s.ProviderById(provId))
}

func (s *stubRepository) BeveragesNeedingSync() []model.Beverage {
	return []model.Beverage{}
}

func (s *stubRepository) SetBeverageMenu(prov model.MenuProvider, beverages []model.Beverage) {
}

func (s *stubRepository) SaveBeverage(beverage model.Beverage) {
}
