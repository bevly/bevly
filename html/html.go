package html

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func Text(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}
	return doc.Text()
}
