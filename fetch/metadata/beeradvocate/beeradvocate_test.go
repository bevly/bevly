package beeradvocate

import (
	"github.com/bevly/bevly/google"
	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func beerAdvocateStub() *httptest.Server {
	return httpfilestub.Server("ba_racer5_test.html")
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
}

func googleSearchStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search, exists := query["q"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if search[0] == "site:beeradvocate.com test" {
			httpfilestub.WriteFile(w, "google_site_search_test.html")
			return
		}
		if search[0] == "beeradvocate test" {
			httpfilestub.WriteFile(w, "google_keyword_search_test.html")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func TestFindProfile(t *testing.T) {
	ts := googleSearchStub()
	defer ts.Close()

	google.SearchBaseURL = ts.URL
	profile, err := FindProfile(model.CreateBeverage("test"))
	assert.Nil(t, err, "FindProfile error")
	assert.Equal(t, "http://www.beeradvocate.com/beer/profile/130/36468/",
		profile, "find cold-hop british")
}
