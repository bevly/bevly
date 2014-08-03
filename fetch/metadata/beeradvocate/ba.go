package metadata

import (
	"github.com/bevly/bevly/google"
	"github.com/bevly/bevly/model"
)

func FetchMetadata(bev model.Beverage) error {
	results, err := google.Search("site:beeradvocate.com " + bev.DisplayName())
	if err != nil {
		return err
	}
	if len(results) > 0 {
		fetchBAMetadata(results[0].URL)
	}
}
