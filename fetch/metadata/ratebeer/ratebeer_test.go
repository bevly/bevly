package ratebeer

import (
	"testing"

	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
)

func TestRatebeerFetch(t *testing.T) {
	ts := httpfilestub.Server("victory_golden_monkey_test.html")
	defer ts.Close()

	bev := model.CreateBeverage("Victory Golden Monkey")
	err := FetchRatebeerMetadata(bev, ts.URL)
	assert.Nil(t, err, "no error from stub")
	assert.True(t, bev.NeedSync(), "should need sync")
	assert.Equal(t, ts.URL, bev.Attribute("rbLink"), "link")
	assert.Equal(t, "Victory Golden Monkey", bev.Name(), "name")
	assert.Equal(t, "Victory Brewing Company", bev.Brewer(), "brewer")
	assert.Equal(t, 97, bev.Ratings()[0].PercentageRating(), "rating")
	assert.Equal(t, 98, bev.Ratings()[1].PercentageRating(), "rating:style")
	assert.Equal(t, 9.5, bev.Abv(), "abv")
	assert.Equal(t, "Enchanting and enlightening, this golden, frothy ale boasts an intriguing herbal aroma, warming alcohol esters on the tongue and light, but firm body to finish. Exotic spices add subtle notes to both the aroma and flavor. Strong, sensual and satisfying.", bev.Description(), "desc")
	assert.Equal(t, "http://res.cloudinary.com/ratebeer/image/upload/w_120,c_limit,q_85,d_no%20image.jpg/beer_630.jpg", bev.Attribute("rbImg"), "image")
}
