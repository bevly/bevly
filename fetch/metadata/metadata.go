package metadata

import "github.com/bevly/bevly/model"

// FetchMetadata fetches metadata for a beverage, from whatever
// sources we deem suitable. Minimal metadata: type of beverage,
// rating.
//
// The beverage object will be modified in-place. If the fetch fails,
// a suitable error object will be returned.
func FetchMetadata(beverage model.Beverage) error {
	return nil
}
