package bing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bevly/bevly/httpfilestub"
)

func bingStub(file string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if len(q["Query"]) == 1 && q["Query"][0] == "'{yak cow}'" {
				httpfilestub.WriteFile(w, file)
				return
			}
			fmt.Fprintf(os.Stderr, "Unexpected query: %v", q["Query"])
			w.WriteHeader(http.StatusBadRequest)
		}))
}

func TestSearch(t *testing.T) {
	ts := bingStub("bing_test.json")
	defer ts.Close()
	s := SearchWithURLKey(ts.URL, "cow")
	res, err := s.Search("yak cow")
	if err != nil {
		t.Errorf("search failed with err: %s", err)
		return
	}

	if len(res) != 10 {
		t.Errorf("expected 10 results, got %d", len(res))
		return
	}
	for _, c := range []struct {
		i    int
		name string
		url  string
	}{
		{0, "Firestone Walker Brewing Company - Union Jack",
			"http://www.firestonebeer.com/beers/products/union-jack"},
		{3, "Firestone Walker Brewing Company - Easy Jack",
			"http://www.firestonebeer.com/beers/products/easy-jack"},
		{9, "Craft Beer: Firestone Walker’s Union Jack IPA Review",
			"http://coed.com/2014/07/25/firestone-walkers-union-jack-ipa-is-whats-on-tap/"},
	} {
		if res[c.i].Text != c.name {
			t.Errorf("result[%d].Text == %s, want %s", c.i, res[c.i].Text,
				c.name)
		}
		if res[c.i].URL != c.url {
			t.Errorf("result[%d].URL == %s, want %s", c.i, res[c.i].URL,
				c.url)
		}
	}
}
