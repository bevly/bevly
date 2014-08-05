package metadata

import (
	"github.com/bevly/bevly/fetch/metadata/beeradvocate"
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/duckduckgo"
	"log"
	"time"
)

// FetchMetadata fetches metadata for a beverage, from whatever
// sources we deem suitable. Minimal metadata: type of beverage,
// rating.
//
// The beverage object will be modified in-place. If the fetch fails,
// a suitable error object will be returned.
func FetchMetadata(beverage model.Beverage) error {
	log.Printf("FetchMetadata: %s", beverage)

	beverage.SetSyncTime(time.Now())
	err := beeradvocate.FetchMetadata(beverage, duckduckgo.DefaultSearch())
	if err != nil && frisco.IsFrisco(beverage) {
		err = frisco.FetchMetadata(beverage)
	}
	return err
}
