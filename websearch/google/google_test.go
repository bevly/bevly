package google

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var racer5Query = "site:beeradvocate.com bear republic racer v"

func googleStubServer(file, search string) *httptest.Server {
	return httpfilestub.ServerValidated(file,
		func(r *http.Request) bool {
			query := r.URL.Query()
			return len(query["q"]) == 1 && query["q"][0] == search
		})
}

func TestSearchWithEmphasizedResult(t *testing.T) {
	ts := googleStubServer("google_greenflash_test.html",
		"beeradvocate Green Flash West Coast IPA")
	defer ts.Close()

	search := SearchWithURL(ts.URL)
	results, err := search.Search("beeradvocate Green Flash West Coast IPA")
	assert.Nil(t, err, "search success")
	assert.Equal(t, 7, len(results), "results")
	assert.Equal(t, "Green Flash West Coast IPA - Beer Advocate",
		results[0].Text, "link text")
	assert.Equal(t, "http://www.beeradvocate.com/beer/profile/2743/22505/",
		results[0].URL, "link url")
}

func TestSearch(t *testing.T) {
	ts := googleStubServer("google_racerv_test.html", racer5Query)
	defer ts.Close()

	search := SearchWithURL(ts.URL)
	results, err := search.Search("site:beeradvocate.com bear republic racer v")
	assert.Nil(t, err, "must search from stub without error")
	assert.Equal(t, 10, len(results), "must find 10 results")
	assert.Equal(t, "Racer 5 India Pale Ale | Bear Republic Brewing Co. - Beer ...", results[0].Text, "must match first result")
	assert.Equal(t, "http://www.beeradvocate.com/beer/profile/610/2751/",
		results[0].URL, "must match first result URL")
	assert.Equal(t, "Black Racer | Bear Republic Brewing Co. - Beer Advocate", results[1].Text, "must match second result")
	assert.Equal(t, "American IPA - Beer Advocate", results[9].Text, "must match last result")
}
