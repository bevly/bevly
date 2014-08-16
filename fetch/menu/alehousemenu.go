package menu

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/text"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const AlehouseDescription = "ale_houseDescription"

func init() {
	menuFetcherRegistry["ale_house"] = alehouseMenu
}

func alehouseMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	agent := httpagent.Agent()
	doc, err := agent.GetDoc(provider.Url())
	if err != nil {
		log.Printf("alehouseMenu: Get(%s) failed: %s\n",
			provider.Url(), err)
		return nil, err
	}

	beverages, err := alehouseDrafts(doc)
	if err != nil {
		log.Printf("alehouseMenu: Failed to parse menu: %s\n", err)
		return nil, err
	}

	log.Printf("alehouseMenu: parsed %d beverages from %s\n",
		len(beverages), provider.Url())
	return beverages, nil
}

type beverageFinder func(sel *goquery.Selection) []model.Beverage

var aleHandlers = map[string]beverageFinder{
	"House Ales":        houseAles,
	"House Ales / Cask": houseCaskAles,
	"Guest Drafts":      guestDrafts,
}

func alehouseDrafts(doc *goquery.Document) ([]model.Beverage, error) {
	ales := doc.Find("section#ales").First().Children()
	if ales.Size() == 0 {
		return nil, ErrEmptyMenu
	}

	res := []model.Beverage{}
	var handler beverageFinder
	ales.Each(func(_ int, sel *goquery.Selection) {
		if sel.Is("h3") {
			sectionHeader := sel.Text()
			handler = aleHandlers[strings.TrimSpace(sectionHeader)]
			if handler != nil {
				log.Printf("alehouseDrafts: Activating handler for %s\n",
					sectionHeader)
			} else {
				log.Printf("alehouseDrafts: No handler for \"%s\"\n",
					sectionHeader)
			}
		} else if handler != nil {
			if sel.Is("div.columns-2") {
				newBevs := handler(sel)
				res = append(res, newBevs...)
			} else {
				html, err := sel.Html()
				if err != nil {
					log.Printf("alehouseDrafts: unexpected non-div.column-2 follows h3: err %s\n", err)
				} else {
					log.Printf("alehouseDrafts: unexpected non-div.column-2 follows h3: %s\n", html)
				}
			}
		}
	})

	if len(res) == 0 {
		return nil, ErrEmptyMenu
	}

	return res, nil
}

func houseAles(sel *goquery.Selection) []model.Beverage {
	return aleMapper(findAles(sel), func(bev model.Beverage) {
		bev.SetDisplayName("Oliver " + bev.DisplayName())
		bev.SetSearchName("Pratt Street Ale House " + bev.DisplayName())
	})
}

func houseCaskAles(sel *goquery.Selection) []model.Beverage {
	return aleMapper(findAles(sel), func(bev model.Beverage) {
		bev.SetDisplayName("Oliver " + bev.DisplayName() + " (cask)")
		bev.SetSearchName("Pratt Street Ale House " + bev.DisplayName())
	})
}

func guestDrafts(sel *goquery.Selection) []model.Beverage {
	return aleMapper(findAles(sel), func(bev model.Beverage) {
		bev.SetDisplayName(strings.Replace(bev.DisplayName(), " - ", " ", 1))
	})
}

func findAles(sel *goquery.Selection) []model.Beverage {
	res := []model.Beverage{}
	sel.Find("li").Each(func(_ int, item *goquery.Selection) {
		title := text.Normalize(item.Find("span").First().Text())
		log.Printf("findAles: Found ale titled: %s\n", title)
		style := text.Normalize(item.Find("p > strong").First().Text())
		desc := text.Normalize(strings.Replace(item.Find("p").First().Text(), style, "", 1))
		abv := findAbv(desc)
		desc = stripStyleAbv(desc)
		bev := model.CreateBeverage(title)
		bev.SetType(style)
		bev.SetDescription(desc)
		bev.SetAttribute(AlehouseDescription, desc)
		if abv > 0.0 {
			bev.SetAbv(abv)
		}
		res = append(res, bev)
	})
	return res
}

var rAlehouseAbv = regexp.MustCompile(`(\d+(?:[.]\d*)?)%\s*$`)

func findAbv(desc string) float64 {
	match := rAlehouseAbv.FindStringSubmatch(desc)
	if match != nil {
		abv, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0.0
		}
		return abv
	}
	return 0.0
}

var rAlehouseDesc = regexp.MustCompile(`(.*?)(\d+(?:[.]\d*)?)%\s*$`)

func stripStyleAbv(desc string) string {
	match := rAlehouseDesc.FindStringSubmatch(desc)
	if match != nil {
		return text.Normalize(match[1])
	} else {
		return text.Normalize(desc)
	}
}

func aleMapper(beverages []model.Beverage, mapper func(model.Beverage)) []model.Beverage {
	for _, bev := range beverages {
		mapper(bev)
	}
	return beverages
}
