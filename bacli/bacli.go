package main

import (
	"fmt"
	"github.com/bevly/bevly/fetch/metadata/beeradvocate"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/websearch/google"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [beer name]\n", os.Args[0])
		os.Exit(1)
	}
	reportBeerMetadata(strings.Join(os.Args[1:], " "))
}

func reportBeerMetadata(beerName string) {
	beer := model.CreateBeverage(beerName)
	beeradvocate.FetchMetadata(beer, google.DefaultSearch())

	fmt.Printf("Title: %s\n", beer.DisplayName())
	fmt.Printf("Name: %s\n", beer.Name())
	fmt.Printf("Brewer: %s\n", beer.Brewer())
	fmt.Printf("Style: %s\n", beer.Type())
	if len(beer.Ratings()) > 0 {
		fmt.Printf("Rating: %s\n", beer.Ratings()[0].PercentageRating())
	}
}
