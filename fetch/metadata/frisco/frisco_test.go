package frisco

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func friscoTestServer() *httptest.Server {
	return httpfilestub.Server("mikkeller_black_test.html")
}

func TestFriscoMetadata(t *testing.T) {
	ts := friscoTestServer()
	defer ts.Close()

	beer := model.CreateBeverage("Mikkeller Black")
	beer.SetAttribute(ProfileURLProperty, ts.URL)
	assert.Nil(t, FetchMetadata(beer), "fetch error")
	assert.Equal(t, "Black 黑", beer.Name(), "name")
	assert.Equal(t, "Mikkeller", beer.Brewer(), "brewer")
	assert.Equal(t, 17.5, beer.Abv(), "ABV")
	assert.Equal(t, "10oz", beer.Attribute(ServingSizeProperty), "serving size")
	assert.Equal(t, "Imperial Stout - The strongest beer in Scandinavia. This Imperial Stout is the craziest, wildest, strongest beer from Mikkeller to this date. Not for sissies.", beer.Description(), "description")
}
