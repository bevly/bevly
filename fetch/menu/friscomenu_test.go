package menu

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 51, len(beverages), "must find all draft beers")
	assert.Equal(t, "Flying Dog Dogtoberfest",
		beverages[7].DisplayName(), "must find Flying Dog Dogtoberfest")
	assert.Equal(t, 5.8, beverages[7].Abv(), "must match Dogtoberfest ABV")
	assert.Equal(t, "16oz", beverages[7].Attribute(frisco.ServingSizeProperty), "must match serving size")
	assert.Equal(t, "Push 72° & Sunny Wheat Ale", beverages[27].DisplayName(),
		"must decode HTML entities")
	assert.Equal(t, ts.URL,
		beverages[50].Attribute(frisco.ProfileURLProperty), "must save detail URL")
}
