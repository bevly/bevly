package frisco

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func friscoTestServer(file string) *httptest.Server {
	return httpfilestub.Server(file)
}

func TestFriscoMetadataMikkeller(t *testing.T) {
	ts := friscoTestServer("mikkeller_black_test.html")
	defer ts.Close()

	beer := model.CreateBeverage("Mikkeller Black")
	beer.SetAttribute(ProfileURLProperty, ts.URL)
	assert.Nil(t, FetchMetadata(beer), "fetch error")
	assert.Equal(t, "Black 黑", beer.Name(), "name")
	assert.Equal(t, "Mikkeller", beer.Brewer(), "brewer")
	assert.Equal(t, 17.5, beer.Abv(), "ABV")
	assert.Equal(t, "10oz", beer.Attribute(ServingSizeProperty), "serving size")
	assert.Equal(t, "Imperial Stout", beer.Type(), "type")
	assert.Equal(t, "The strongest beer in Scandinavia. This Imperial Stout is the craziest, wildest, strongest beer from Mikkeller to this date. Not for sissies.", beer.Description(), "description")
}

func TestFriscoMetadataBoulevard(t *testing.T) {
	ts := friscoTestServer("boulevard_imperial_test.html")
	defer ts.Close()
	beer := model.CreateBeverage("Boulevard Imperial Stout")
	beer.SetAttribute(ProfileURLProperty, ts.URL)
	assert.Nil(t, FetchMetadata(beer), "fetch error")
	assert.Equal(t, "Imperial Stout", beer.Name(), "name")
	assert.Equal(t, "Boulevard", beer.Brewer(), "brewer")
	assert.Equal(t, 11.8, beer.Abv(), "ABV")
	assert.Equal(t, "10oz", beer.Attribute(ServingSizeProperty), "serving")
	assert.Equal(t, "Imperial Stout", beer.Type(), "type")
	assert.Equal(t, "63", beer.Attribute(IBUProperty), "ibu")
	assert.Equal(t, "Like India Pale Ale, Imperial Stouts were originally brewed with high percentages of malt and hops to help them withstand the rigors of a long sea journey, not to India but to Imperial Russia and the Baltic States. This version is an over-the-top riff on the style, with a huge grain bill featuring several kinds of malted barley, wheat, rye, and oats. A portion of the batch has been aged in whiskey barrels, then blended with fresh ale before bottling/kegging.", beer.Description(), "desc")
}
