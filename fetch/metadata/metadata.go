package metadata

import (
	"github.com/bevly/bevly/fetch/metadata/beeradvocate"
	"github.com/bevly/bevly/fetch/metadata/frisco"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/google"
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
	err := beeradvocate.FetchMetadata(beverage, google.DefaultSearch())
	if frisco.IsFrisco(beverage) {
		friscoErr := frisco.FetchMetadata(beverage)
		if friscoErr != nil {
			// Frisco fetch errors are harmless, just log for
			// investigation:
			log.Printf("FetchMetadata(%s): frisco fetch error: %s\n",
				beverage, friscoErr)
		}
	}
	return err
}
