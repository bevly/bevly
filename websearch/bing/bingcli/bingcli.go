package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bevly/bevly/websearch/bing"
)

func main() {
	search := strings.Join(os.Args[1:], " ")
	if search == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [search terms]", os.Args[0])
		os.Exit(1)
	}
	s := bing.DefaultSearch()
	res, err := s.Search(search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search for %s failed: %s\n", search, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Results: %v\n", res)
}
