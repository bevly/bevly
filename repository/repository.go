package repository

import (
	"log"

	"github.com/bevly/bevly/model"
)

type Repository interface {
	MenuProviders() []model.MenuProvider
	ProviderByID(id string) model.MenuProvider
	ProviderBeverages(provider model.MenuProvider) []model.Beverage
	ProviderIDBeverages(providerName string) []model.Beverage
	BeveragesNeedingSync() []model.Beverage
	BeverageByName(name string) model.Beverage

	// TODO
	SetBeverageMenu(provider model.MenuProvider, menu []model.Beverage)
	SaveBeverage(beverage model.Beverage)
	// BeverageByName(name string)

	// Discard unreferenced beverages
	GarbageCollect()

	// Delete everything in the repository
	Purge()
}

func StubRepository() Repository {
	return &stubRepository{}
}

type stubRepository struct{}

var _ Repository = &stubRepository{}

func (s *stubRepository) MenuProviders() []model.MenuProvider {
	return []model.MenuProvider{
		model.CreateMenuProvider("frisco", "Frisco", "http://beer.friscogrille.com/", "frisco"),
		model.CreateMenuProvider("ale_house", "Ale House", "http://www.thealehousecolumbia.com/menu/", "ale_house"),
	}
}

func (s *stubRepository) ProviderByID(id string) model.MenuProvider {
	log.Printf("Looking for provider named \"%s\"\n", id)
	for _, prov := range s.MenuProviders() {
		if prov.ID() == id {
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

func (s *stubRepository) ProviderIDBeverages(provID string) []model.Beverage {
	return s.ProviderBeverages(s.ProviderByID(provID))
}

func (s *stubRepository) BeveragesNeedingSync() []model.Beverage {
	return []model.Beverage{}
}

func (s *stubRepository) SetBeverageMenu(prov model.MenuProvider, beverages []model.Beverage) {
}

func (s *stubRepository) SaveBeverage(beverage model.Beverage) {
}

func (*stubRepository) GarbageCollect() {
}

func (*stubRepository) Purge() {
}

func (*stubRepository) BeverageByName(name string) model.Beverage {
	return nil
}
