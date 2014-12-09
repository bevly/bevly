package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bevly/bevly/websearch/duckduckgo"
)

func main() {
	search := strings.Join(os.Args[1:], " ")
	res, err := duckduckgo.DefaultSearch().Search(search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search for %s failed: %s\n", search, err)
		os.Exit(1)
	}
	fmt.Println("Found", len(res), "results for", search)
	for i, r := range res {
		fmt.Printf("%d) %s (%s)\n", i, r.Text, r.URL)
	}
}
