package repository

import "model"

type Repository interface {
	MenuProviders() []model.MenuProvider
	BeverageMenu(provider model.MenuProvider) []model.Beverage

	// TODO
	// SetBeverageMenu(provider model.MenuProvider, menu []model.Beverage)
	// SetBeverage(beverage model.Beverage)
	// BeverageByName(name string)

	// Clear()
}

func DefaultRepository() Repository {
	return &stubRepository{}
}

type stubRepository struct{}

func (s *stubRepository) MenuProviders() []model.MenuProvider {
	return []model.MenuProvider{
		model.CreateMenuProvider("Frisco", "http://beer.friscogrille.com/",
			"frisco"),
	}
}

func (s *stubRepository) BeverageMenu() []model.Beverage {
	return []model.Beverage{
		model.CreateBeverageAbvTypeRatingLink(
			"Mikkeller Beer Geek Brunch Weasel",
			10.9,
			"", 0, "",
			"http://www.google.com/search?q=Mikkeller+Beer+Geek+Brunch+Weasel"),

		model.CreateBeverageAbvTypeRatingLink(
			"Dogfish Head Olde School Barleywine",
			15.0,
			"", 0, "",
			"http://www.google.com/search?q=Dogfish+Head+Olde+School+Barleywine"),

		model.CreateBeverageAbvTypeRatingLink(
			"Blue Point Toasted Lager",
			5.3,
			"", 80, "BA", "http://beeradvocate.com/beer/profile/764/2318/"),
	}
}
