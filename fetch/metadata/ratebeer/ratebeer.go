package ratebeer

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"github.com/bevly/bevly/websearch"
	"strconv"
)

var ErrNoResults = errors.New("no results for beverage")

const RateBeerAccuracyScore = 9

func FetchMetadata(bev model.Beverage, search websearch.Search) (err error) {
	log.Printf("FetchMetadata(%s): Searching for Ratebeer profile", bev)
	profileURL, err := FindProfile(bev.DisplayName(), search)
	if err != nil {
		log.Printf("FetchMetadata(%s): Ratebeer profile error: %s",
			bev, err)
		return
	}
	if profileURL == "" {
		return ErrNoResults
	}
	return FetchRatebeerMetadata(bev, profileURL)
}

func FetchRatebeerMetadata(bev model.Beverage, profileURL string) (err error) {
	doc, err := httpagent.Agent().GetDoc(profileURL)
	if err != nil {
		return
	}
	bev.SetNeedSync(true)

	overwrite := bev.AccuracyScore() < RateBeerAccuracyScore
	if overwrite {
		bev.SetAccuracyScore(RateBeerAccuracyScore)
	}

	set := func(old, new string, setter func(string)) {
		if new != "" && (overwrite || old == "") {
			setter(new)
		}
	}

	selFirstText := func(selector string) string {
		return doc.Find(selector).First().Text()
	}

	bev.SetAttribute("rbLink", profileURL)
	set(bev.Link(), profileURL, bev.SetLink)
	set(bev.Name(), selFirstText(".user-header h1"), bev.SetName)
	set(bev.Brewer(), selFirstText("big a"), bev.SetBrewer)
	addRatings(bev, doc)

	abv := findAbv(doc)
	if abv > 0.0 && (overwrite || bev.Abv() == 0.0) {
		bev.SetAbv(abv)
	}

	desc := findDescription(doc)
	set(bev.Description(), desc, bev.SetDescription)
	bev.SetAttribute("rbDescription", desc)

	image := findImageUrl(doc)
	if image != "" {
		bev.SetAttribute("rbImg", image)
		set(bev.Attribute("img"), image, func(img string) {
			bev.SetAttribute("img", img)
		})
	}

	return nil
}

func findImageUrl(doc *goquery.Document) string {
	container := doc.Find("a[href=\"/PictureCredits.asp\"]").Parent().Parent()
	img, exists := container.Find("img").First().Attr("src")
	if exists {
		return img
	}
	return ""
}

const DescriptionHeading = "COMMERCIAL DESCRIPTION"

func findDescription(doc *goquery.Document) string {
	desc := ""
	doc.Find("div > small").EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		if sel.Text() == DescriptionHeading {
			blockText := sel.Parent().Text()
			desc = text.Normalize(
				strings.Replace(blockText, DescriptionHeading, "", 1))
			return false
		}
		return true
	})
	return desc
}

var rAbv = regexp.MustCompile(`(\d+(?:[.]\d+)?)%`)

func findAbv(doc *goquery.Document) float64 {
	abvText := doc.Find("abbr[title=\"Alcohol By Volume\"]").Next().Text()
	abvMatch := rAbv.FindStringSubmatch(abvText)
	if abvMatch != nil {
		abv, err := strconv.ParseFloat(abvMatch[1], 64)
		if err != nil || abv <= 0.0 {
			return 0.0
		}
		return abv
	}
	return 0.0
}

func addRatings(bev model.Beverage, doc *goquery.Document) {
	ratings := doc.Find("[itemtype=\"http://data-vocabulary.org/Rating\"]").First().Children()
	log.Printf("rating text: %s\n", ratings.Text())

	parseNum := func(word, desc string) int {
		log.Printf("Parsing rating from %s (- %s)\n", desc, word)
		stripped := text.Normalize(strings.Replace(desc, word, "", 1))
		if stripped != "" {
			res, err := strconv.ParseInt(stripped, 10, 32)
			if err != nil {
				return 0
			}
			return int(res)
		}
		return 0
	}

	overall := parseNum("overall", ratings.First().Text())
	style := parseNum("style", ratings.First().Next().Text())
	if overall > 0 {
		bev.AddRating(model.CreateRating("rb", overall))
	}
	if style > 0 {
		bev.AddRating(model.CreateRating("rb:style", style))
	}
}

func FindProfile(name string, search websearch.Search) (string, error) {
	terms := "ratebeer " + name
	results, err := search.Search(terms)
	if err != nil {
		return "", err
	}
	log.Printf("ratebeer(%s): %d results\n", terms, len(results))
	if len(results) > 0 {
		for _, result := range results {
			log.Printf("rb(%s): considering %s (%s)\n",
				search, result.Text, result.URL)
			if IsRatebeerProfile(result.URL) {
				title := stripRatebeerName(result.Text)
				confidence := text.NameMatchConfidence(title, name)
				if confidence < 0.13 {
					log.Printf("rb(%s): rejecting %s (confidence: %.2f%%)\n",
						terms, title, confidence*100)
					continue
				}
				log.Printf("rb(%s): accepting %s (confidence: %.2f%%)\n",
					terms, title, confidence*100)
				return result.URL, nil
			}
		}
	}
	return "", nil
}

var rRatebeerURL = regexp.MustCompile(`ratebeer.com/beer/.*?/\d+/?$`)

func IsRatebeerProfile(url string) bool {
	return rRatebeerURL.FindString(url) != ""
}

func stripRatebeerName(text string) string {
	return strings.Replace(text, "- RateBeer", "", 1)
}
