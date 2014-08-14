package http

import (
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/repository"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"net/http"
)

func BeverageServerBlocking(repo repository.Repository) {
	m := martini.Classic()
	m.Use(gzip.All())
	m.Use(render.Renderer())

	m.Get("/:source/drink/", func(par martini.Params, r render.Render, res http.ResponseWriter) {
		NoCache(res)
		r.JSON(http.StatusOK, bevListJsonModel(repo.ProviderIdBeverages(par["source"])))
	})
	m.Run()
}

func NoCache(res http.ResponseWriter) {
	headers := res.Header()
	headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	headers.Set("Pragma", "no-cache")
	headers.Set("Expires", "0")
}

func bevListJsonModel(beverages []model.Beverage) interface{} {
	bevList := make([]interface{}, len(beverages))
	for i, beverage := range beverages {
		bevList[i] = bevJsonModel(beverage)
	}
	return map[string]interface{}{
		"drinks": bevList,
	}
}

func simpleRatingScore(beverage model.Beverage) interface{} {
	for _, rating := range beverage.Ratings() {
		return rating.PercentageRating()
	}
	return nil
}

func bevJsonModel(beverage model.Beverage) interface{} {
	bevJson := map[string]interface{}{
		"name":         beverage.DisplayName(),
		"brewer":       beverage.Brewer(),
		"type":         beverage.Type(),
		"abv":          beverage.Abv(),
		"description":  beverage.Description(),
		"externalLink": beverage.Link(),
		"ratingScore":  simpleRatingScore(beverage),
	}

	// Attributes should be named to not collide:
	for name, value := range beverage.Attributes() {
		bevJson[name] = value
	}

	return bevJson
}
