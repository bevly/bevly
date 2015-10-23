package mongorepo

import (
	"testing"

	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/repository"
	"github.com/stretchr/testify/assert"
)

var repo repository.Repository

func init() {
	var err error
	repo, err = Repository("127.0.0.1", "bevtest")
	if err != nil {
		panic(err)
	}
}

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
	bev := beverageInfos[0].Model()
	bev.SetLink("http://foo")
	repo.SaveBeverage(bev)
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
	bevModel.SetType("cow")
	repo.SaveBeverage(bevModel)

	bevModel = beverageInfos[0].Model()
	bevModel.SetAbv(0.0)
	bevModel.SetDescription("Hii")
	bevModel.SetAttribute("cow", "moo")
	bevModel.AddRating(model.CreateRating("BA", 95))
	bevModel.AddRating(model.CreateRating("RateBeer", 87))
	bevModel.SetLink("http://google.com")
	bevModel.SetType("IPA")
	bevModel.SetAccuracyScore(10)
	repo.SaveBeverage(bevModel)

	bev := repo.BeverageByName(bevModel.DisplayName())
	if assert.NotNil(t, bev, "Updated beverage should be discovered") {
		assert.Equal(t, "IPA", bev.Type(), "type")
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
		assert.Equal(t, "http://google.com", bev.Link(), "link")
	}

	bevModel.SetAccuracyScore(0)
	bevModel.SetAbv(1.1)
	bevModel.SetType("cowboy")
	repo.SaveBeverage(bevModel)
	bev = repo.BeverageByName(bevModel.DisplayName())
	assert.Equal(t, 10, bev.AccuracyScore(), "score")
	assert.Equal(t, 4.54, bev.Abv(), "ABV preserve")
	assert.Equal(t, "IPA", bev.Type(), "type preserve")
}

func TestSaveMenu(t *testing.T) {
	repo.Purge()
	frisco := repo.ProviderByID("frisco")
	bevs := make([]model.Beverage, len(beverageInfos))
	for i, bevInfo := range beverageInfos {
		bevs[i] = bevInfo.Model()
	}
	repo.SetBeverageMenu(frisco, bevs)

	savedBevs := repo.ProviderIDBeverages("frisco")
	assert.Equal(t, 3, len(savedBevs), "three beverages should be saved")
	assert.Equal(t, "Bear Republic Racer V", savedBevs[1].DisplayName())
}
