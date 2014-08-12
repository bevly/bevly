package menu

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func alehouseStub() *httptest.Server {
	return httpfilestub.Server("alehouse_test.html")
}

func TestAlehouseMenu(t *testing.T) {
	ts := alehouseStub()

	prov := model.CreateMenuProvider(
		"ale_house", "Ale House", ts.URL, "ale_house")
	bev, err := alehouseMenu(prov)
	assert.Nil(t, err, "stub fetch must succeed")
	assert.Equal(t, 30, len(bev), "must find all beers")
	assert.Equal(t, "Oliver “The Black Code”", bev[0].DisplayName(), "black code")
	assert.Equal(t, "Oliver King of the Knight Time World", bev[2].DisplayName(), "knight time")
	assert.Equal(t, "Oliver “The Black Code” (cask)", bev[15].DisplayName(), "black code cask")
	assert.Equal(t, "Allagash White Ale", bev[29].DisplayName(), "allagash white")

	assert.Equal(t, "Black Ale", bev[0].Type(), "black code type")
	assert.Equal(t, "Belgian Style White Ale", bev[29].Type(), "allagash white type")

	assert.Equal(t, 6.5, bev[0].Abv(), "black code ABV")
	assert.Equal(t, 5, bev[29].Abv(), "allagash white ABV")
	assert.Equal(t, "Brewed with a generous portion of wheat and spiced with coriander and Curaco orange peel, this beer is refreshing and slightly cloudy in appearance.", bev[29].Description(), "allagash white desc")

	assert.Equal(t, 7.7, bev[1].Abv(), "sentinel ABV")
	assert.Equal(t, 8.4, bev[2].Abv(), "knight ABV")
}
