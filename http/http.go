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
	m.Use(render.Renderer())
	m.Use(gzip.All())

	m.Get("/:source/drink/", func(par martini.Params, r render.Render) {
		r.JSON(http.StatusOK, bevListJsonModel(repo.ProviderIdBeverages(par["source"])))
	})
	m.Run()
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
	return map[string]interface{}{
		"name":         beverage.DisplayName(),
		"brewer":       beverage.Brewer(),
		"type":         beverage.Type(),
		"abv":          beverage.Abv(),
		"externalLink": beverage.Link(),
		"ratingScore":  simpleRatingScore(beverage),
	}
}
