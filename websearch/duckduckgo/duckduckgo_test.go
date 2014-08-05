package duckduckgo

import (
	"github.com/bevly/bevly/httpfilestub"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func duckStub() *httptest.Server {
	return httpfilestub.Server("cider_test.js")
}

func TestSearch(t *testing.T) {
	ts := duckStub()
	defer ts.Close()

	search := SearchWithURL(ts.URL)
	results, err := search.Search("cider")
	assert.Nil(t, err, "must search from stub without error")
	assert.Equal(t, 26, len(results), "26 results")
	assert.Equal(t, "Our Hard Ciders | Bold Rock Hard Cider",
		results[0].Text, "first result title")
	assert.Equal(t, "http://www.boldrock.com/OurCiders.html",
		results[0].URL, "first result URL")

	assert.Equal(t, "1K Beer Walk | Crystal City BID | Arlington, VA",
		results[25].Text, "last result title")
	assert.Equal(t, "http://www.crystalcity.org/do/1k-beer-walk",
		results[25].URL, "last result URL")
}
