package metadata

import (
	"log"
	"time"

	"github.com/bevly/bevly/fetch/metadata/beeradvocate"
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/fetch/metadata/ratebeer"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/duckduckgo"
)

// FetchMetadata fetches metadata for a beverage, from whatever
// sources we deem suitable. Minimal metadata: type of beverage,
// rating.
//
// The beverage object will be modified in-place. If the fetch fails,
// a suitable error object will be returned.
func FetchMetadata(beverage model.Beverage) (err error) {
	log.Printf("FetchMetadata: %s", beverage)

	beverage.SetSyncTime(time.Now())

	search := duckduckgo.DefaultSearch()
	err = ratebeer.FetchMetadata(beverage, search)
	if err != nil {
		log.Printf("FetchMetadata(%s): ratebeer fetch error: %s\n",
			beverage, err)
	}

	if frisco.IsFrisco(beverage) {
		err = frisco.FetchMetadata(beverage)
		if err != nil {
			log.Printf("FetchMetadata(%s): frisco fetch error: %s\n",
				beverage, err)
		}
	}

	err = beeradvocate.FetchMetadata(beverage, search)
	return err
}
