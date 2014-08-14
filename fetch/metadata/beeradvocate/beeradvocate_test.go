package beeradvocate

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/duckduckgo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func beerAdvocateStub() *httptest.Server {
	return httpfilestub.Server("ba_racer5_test.html")
}

func TestIsBeerAdvocateProfile(t *testing.T) {
	assert.True(t, IsBeerAdvocateProfile("http://www.beeradvocate.com/beer/profile/2743/34085/"), "is ba profile")
	assert.False(t, IsBeerAdvocateProfile("http://www.beeradvocate.com/beer/profile/2743/"), "is NOT ba profile")
}

func TestBAMetadata(t *testing.T) {
	ts := beerAdvocateStub()
	defer ts.Close()

	beer := model.CreateBeverage("Bear Republic Racer V")
	assert.Nil(t, fetchBAMetadata(beer, ts.URL), "metadata fetch from stub must not fail")
	assert.Equal(t, ts.URL, beer.Link(), "link == metadata url")
	assert.Equal(t, "Bear Republic Brewing Co.", beer.Brewer(), "brewer")
	assert.Equal(t, "Racer 5 India Pale Ale", beer.Name(), "name")
	assert.Equal(t, 7.5, beer.Abv(), "abv")
	assert.Equal(t, "American IPA", beer.Type(), "style")
	assert.Equal(t, 94, beer.Ratings()[0].PercentageRating(), "BA rating")
	assert.Equal(t, 93, beer.Ratings()[1].PercentageRating(), "BA bro rating")
	assert.Equal(t, "BAbro", beer.Ratings()[1].Source(), "BA bro rating name")
	assert.Equal(t, "This hoppy American IPA is a full bodied beer brewed with American pale and crystal malts, and heavily hopped with Chinook, Cascade, Columbus and Centennial. There's a trophy in every glass. \n \n2009 Great American Beer Festival® American-Style Strong Pale Ale – GOLD \n2009 Colorado State Fair – Best of Show \n2006 Great American Beer Festival® American-Style Strong Pale Ale – SILVER \n2005 Great American Beer Festival® American-Style Strong Pale Ale – SILVER \n1999 Great American Beer Festival® India Pale Ale - GOLD", beer.Description(), "desc")
}

func webSearchStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search, exists := query["q"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if search[0] == "beeradvocate boulder cold hop english ipa" {
			httpfilestub.WriteFile(w, "boulder_duck.js")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func TestFindProfile(t *testing.T) {
	ts := webSearchStub()
	defer ts.Close()

	search := duckduckgo.SearchWithURL(ts.URL)
	profile, err := FindProfile(model.CreateBeverage("boulder cold hop english ipa"), search)
	assert.Nil(t, err, "FindProfile error")
	assert.Equal(t, "http://www.beeradvocate.com/beer/profile/130/36468/",
		profile, "find cold-hop british")
}
