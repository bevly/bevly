package ratebeer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	Throttle.SetDisabled(true)
}

func redirectStub(file string, redirectPath string, redirectFile string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == redirectPath {
			httpfilestub.WriteFile(w, redirectFile)
		} else {
			httpfilestub.WriteFile(w, file)
		}
	}))
}

func TestRatebeerRedirectFetch(t *testing.T) {
	ts := redirectStub("woodchuck-hopped-redirect_test.html",
		"/beer/woodchuck-cellar-series-01--dry-hop/280087/",
		"woodchuck-hopsation_test.html")
	defer ts.Close()

	bev := model.CreateBeverage("Woodchuck Hopped Apple")
	err := FetchRatebeerMetadata(bev, ts.URL)
	assert.Nil(t, err, "no error from stub")
	assert.True(t, bev.NeedSync(), "should need sync")
	assert.Equal(t, "Woodchuck Hopsation", bev.Name(), "name")
	assert.Equal(t, "Vermont Hard Cider Company, LLC (C&C Group)",
		bev.Brewer(), "brewer")
	assert.Equal(t, 6.9, bev.Abv(), "abv")
}

func TestRatebeerFetch(t *testing.T) {
	ts := httpfilestub.Server("victory_golden_monkey_test.html")
	defer ts.Close()

	bev := model.CreateBeverage("Victory Golden Monkey")
	err := FetchRatebeerMetadata(bev, ts.URL)
	assert.Nil(t, err, "no error from stub")
	assert.True(t, bev.NeedSync(), "should need sync")
	assert.Equal(t, ts.URL, bev.Attribute("rbLink"), "link")
	assert.Equal(t, "Victory Golden Monkey", bev.Name(), "name")
	assert.Equal(t, "Victory Brewing Company", bev.Brewer(), "brewer")
	assert.Equal(t, 97, bev.Ratings()[0].PercentageRating(), "rating")
	assert.Equal(t, 98, bev.Ratings()[1].PercentageRating(), "rating:style")
	assert.Equal(t, 9.5, bev.Abv(), "abv")
	assert.Equal(t, "Enchanting and enlightening, this golden, frothy ale boasts an intriguing herbal aroma, warming alcohol esters on the tongue and light, but firm body to finish. Exotic spices add subtle notes to both the aroma and flavor. Strong, sensual and satisfying.", bev.Description(), "desc")
	assert.Equal(t, "http://res.cloudinary.com/ratebeer/image/upload/w_120,c_limit,q_85,d_no%20image.jpg/beer_630.jpg", bev.Attribute("rbImg"), "image")
}

func beverage(beerName, file string) (model.Beverage, error) {
	ts := httpfilestub.Server(file)
	defer ts.Close()
	bev := model.CreateBeverage(beerName)
	err := FetchRatebeerMetadata(bev, ts.URL)
	return bev, err
}

func TestRatebeerRating(t *testing.T) {
	name := "Abita Bourbon Street Maple Pecan"
	bev, err := beverage(name, "abita-maple-pecan_test.html")
	if err != nil {
		t.Errorf("error creating beverage: %#v: %s", name, err)
		return
	}
	if len(bev.Ratings()) != 2 {
		t.Errorf("Abita Bourbon: expected 2 ratings (overall and style), got: %d", len(bev.Ratings()))
		return
	}
	if bev.Ratings()[0].PercentageRating() != 75 {
		t.Errorf("Abita Bourbon: expected overall rating == 75, got %d", bev.Ratings()[0].PercentageRating())
	}
	if bev.Ratings()[1].PercentageRating() != 85 {
		t.Errorf("Abita Bourbon: expected overall rating == 85, got %d", bev.Ratings()[1].PercentageRating())
	}
}

func TestRatebeerEncoding(t *testing.T) {
	ts := httpfilestub.Server("southern_tier_pumking_test.html")
	defer ts.Close()

	bev := model.CreateBeverage("Southern Tier Pumking")
	err := FetchRatebeerMetadata(bev, ts.URL)
	assert.Nil(t, err, "no error from stub")
	assert.Equal(t, "Pumking is an ode to Púca, a creature of Celtic folklore, who is both feared and respected by those who believe in it. Púca is said to waylay travelers throughout the night, tossing them on its back, and providing them the ride of their lives, from which they return forever changed. Brewed in the spirit of All Hallows Eve, a time of the year when spirits can make contact with the physical world and when magic is most potent. Pour Pumking into a goblet and allow it’s alluring spirit to overflow. As spicy aromas present themselves, let it’s deep copper color entrance you as your journey into this mystical brew has just begun. As the first drops touch your tongue a magical spell will bewitch your taste buds making it difficult to escape the Pumking. 2007 - Brown Label 7.9% ABV w/text & logo on bottlecap 2008 - Brown Label 9.0% ABV w/text & logo on bottlecap 2009 - Orange Label 9.0% ABV w/text & logo on bottlecap 2010 - Orange Label 9.0% ABV w/logo only on bottlecap 2011 - Orange wood grain background label 8.6% ABV and Southern Tier logotype in two lines. 2012 - Same as 2011 label 8.6% ABV, silver/black/white bottlecap. 2013 - Same as 2012 label and bottle cap 8.6% ABV, date stamp in green text. 2014 - New label design w/orange/white/green, 8.6% ABV", bev.Description(), "description")
}
