package beeradvocate

import (
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func beerAdvocateStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile("ba_racer5_test.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(bytes)
		}
	}))
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
