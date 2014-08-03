package mongorepo

import (
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var repo = DefaultRepository()

type beverageInfo struct {
	name         string
	abv          float64
	bevType      string
	rating       int
	ratingSource string
	link         string
}

func (b *beverageInfo) Model() model.Beverage {
	return model.CreateBeverageAbvTypeRatingLink(b.name, b.abv, b.bevType, b.rating, b.ratingSource, b.link)
}

var beverageInfos = []beverageInfo{
	{"Anchor IPA", 4.54, "IPA", 90, "BA", "http://cow.org"},
	{"Bear Republic Racer V", 4.7, "IPA", 95, "BA", "http://ba.org"},
	{"Jolly Pumpkin Oro de Calabaza", 6.0, "Bière de Garde", 92, "BA",
		"http://www.beeradvocate.com/beer/profile/9897/18975/"},
}

func saveDefaultBeverage() {
	repo.SaveBeverage(beverageInfos[0].Model())
}

func TestSaveBeverage(t *testing.T) {
	repo.Purge()
	saveDefaultBeverage()
	bev := repo.BeverageByName(beverageInfos[0].name)
	if assert.NotNil(t, bev, "Saved beverage should be discovered") {
		assert.Equal(t, "Anchor IPA", bev.DisplayName(), "saved name incorrect")
	}
}

func TestSaveBeverageUpdate(t *testing.T) {
	repo.Purge()
	saveDefaultBeverage()

	bevModel := beverageInfos[0].Model()
	bevModel.SetAbv(0.0)
	bevModel.SetDescription("Hii")
	bevModel.SetAttribute("cow", "moo")
	bevModel.AddRating(model.CreateRating("BA", 95))
	bevModel.AddRating(model.CreateRating("RateBeer", 87))
	repo.SaveBeverage(bevModel)

	bev := repo.BeverageByName(bevModel.DisplayName())
	if assert.NotNil(t, bev, "Updated beverage should be discovered") {
		assert.Equal(t, 4.54, bev.Abv(), "ABV should be unmodified")
		assert.Equal(t, 2, len(bev.Ratings()), "Rating count should be 2")
		assert.Equal(t, 95, bev.Ratings()[0].PercentageRating(),
			"Rating should be updated in place")
		assert.Equal(t, 87, bev.Ratings()[1].PercentageRating(),
			"Ratebeer rating should be saved")
		assert.Equal(t, "RateBeer", bev.Ratings()[1].Source(),
			"Ratebeer source should be saved")
		assert.Equal(t, "Hii", bev.Description(), "description")
		assert.Equal(t, "moo", bev.Attribute("cow"), "attribute:cow")
	}
}

func TestSaveMenu(t *testing.T) {
	repo.Purge()
	frisco := repo.ProviderById("frisco")
	bevs := make([]model.Beverage, len(beverageInfos))
	for i, bevInfo := range beverageInfos {
		bevs[i] = bevInfo.Model()
	}
	repo.SetBeverageMenu(frisco, bevs)

	savedBevs := repo.ProviderIdBeverages("frisco")
	assert.Equal(t, 3, len(savedBevs), "three beverages should be saved")
	assert.Equal(t, "Bear Republic Racer V", savedBevs[1].DisplayName())
}
