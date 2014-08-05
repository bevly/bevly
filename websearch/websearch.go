package websearch

type Result struct {
	URL  string
	Text string
}

type Search interface {
	Search(terms string) ([]Result, error)
	SearchURL(terms string) string
}
