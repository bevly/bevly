package google

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var racer5Query = "site:beeradvocate.com bear republic racer v"

func googleStubServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if len(query["q"]) != 1 || query["q"][0] != racer5Query {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, err := ioutil.ReadFile("google_racerv_test.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(bytes)
		}
	}))
}

func TestSearch(t *testing.T) {
	ts := googleStubServer()
	defer ts.Close()

	SearchBaseURL = ts.URL
	results, err := Search("site:beeradvocate.com bear republic racer v")
	assert.Nil(t, err, "must search from stub without error")
	assert.Equal(t, 10, len(results), "must find 10 results")
	assert.Equal(t, "Racer 5 India Pale Ale | Bear Republic Brewing Co. - Beer ...", results[0].Text, "must match first result")
	assert.Equal(t, "http://www.beeradvocate.com/beer/profile/610/2751/",
		results[0].URL.String(), "must match first result URL")
	assert.Equal(t, "Black Racer | Bear Republic Brewing Co. - Beer Advocate", results[1].Text, "must match second result")
	assert.Equal(t, "American IPA - Beer Advocate", results[9].Text, "must match last result")
}
