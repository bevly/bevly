package menu

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/bevly/bevly/httpfilestub"
	"github.com/bevly/bevly/model"
	"github.com/stretchr/testify/assert"
)

func alehouseStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/menu/" {
			rehost := "http://" + r.Host + "/"
			fmt.Fprintf(os.Stderr, "Serving alehouse_test.html, rehosted as: %s\n", rehost)
			httpfilestub.WriteFileMunged(w, "alehouse_test.html",
				regexp.MustCompile(`http://www[.]thealehousecolumbia[.]com/`),
				rehost)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".pdf") {
			httpfilestub.WriteFile(w, "alemenu_test.pdf")
			return
		}
	}))
}

func TestAlehouseMenu(t *testing.T) {
	ts := alehouseStub()

	prov := model.CreateMenuProvider(
		"ale_house", "Ale House", ts.URL+"/menu/", "ale_house")
	bev, err := alehouseMenu(prov)
	assert.Nil(t, err, "stub fetch must succeed")
	assert.Equal(t, 32, len(bev), "must find all beers")
	assert.Equal(t, "Jailbreak Desserted", bev[0].DisplayName())
	assert.Equal(t, "Jailbreak Feed The Monkey", bev[2].DisplayName())
	assert.Equal(t, "Lagunitas Undercover Investigation Shut-Down", bev[15].DisplayName())
	assert.Equal(t, "Oliver 3 Lions", bev[29].DisplayName())

	assert.Equal(t, "Chocolate Coconut Porter", bev[0].Type())
	assert.Equal(t, "Strong Brown Ale", bev[29].Type())

	assert.Equal(t, 6.9, bev[0].Abv())
	assert.Equal(t, 7.5, bev[29].Abv())
	assert.Equal(t, "A strong English style brown ale, full-bodied with underlying hints of caramel sweetness.",
		bev[29].Description())

	assert.Equal(t, 8.5, bev[1].Abv())
	assert.Equal(t, 6, bev[2].Abv())
}
