package menu

import (
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func friscoStubServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile("frisco_test.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(bytes)
		}
	}))
}

func TestFriscoMenu(t *testing.T) {
	ts := friscoStubServer()
	defer ts.Close()

	provider := model.CreateMenuProvider("frisco", "Frisco", ts.URL, "frisco")
	beverages, err := FetchMenu(provider)
	assert.Nil(t, err, "fetch from stub server must not fail")
	assert.Equal(t, 56, len(beverages), "must find all draft beers")
	assert.Equal(t, "Dieu Du Ciel Péché Mortel Imperial Stout",
		beverages[7].DisplayName(), "must find Dieu Du Ciel")
	assert.Equal(t, "Mikkeller Black 黑", beverages[29].DisplayName(),
		"must decode HTML entities")
	assert.Equal(t, "Union Double Duckpin", beverages[54].DisplayName(),
		"must strip embedded double spaces")
	assert.Equal(t, ts.URL+"/beer_details.php?beer_id=3468",
		beverages[54].Attribute(frisco.ProfileURLProperty), "must save detail URL")
}
