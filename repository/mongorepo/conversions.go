package mongorepo

import (
	"github.com/bevly/bevly/model"
)

func repoBeverageModels(repoBevs []repoBeverage) []model.Beverage {
	result := make([]model.Beverage, len(repoBevs))
	for i, repoBev := range repoBevs {
		result[i] = repoBeverageModel(&repoBev)
	}
	return result
}

func repoBeverageModel(repoBev *repoBeverage) model.Beverage {
	bev := model.CreateBeverage(repoBev.DisplayName)
	bev.SetType(repoBev.BevType)
	bev.SetBrewer(repoBev.Brewer)
	bev.SetLink(repoBev.Link)
	bev.SetAbv(repoBev.Abv)
	bev.SetAttributes(repoBev.Attributes)
	for _, rating := range repoBev.Ratings {
		bev.AddRating(model.CreateRating(rating.Source, rating.PercentageRating))
	}
	return bev
}

func beverageModelToRepo(bev model.Beverage) *repoBeverage {
	repoBev := &repoBeverage{}
	repoBev.DisplayName = bev.DisplayName()
	repoBev.BevType = bev.Type()
	repoBev.Brewer = bev.Brewer()
	repoBev.Link = bev.Link()
	repoBev.Abv = bev.Abv()
	repoBev.Attributes = bev.Attributes()

	for _, rating := range bev.Ratings() {
		repoBev.Ratings = append(repoBev.Ratings,
			repoRating{
				Source:           rating.Source(),
				PercentageRating: rating.PercentageRating(),
			})
	}
	return repoBev
}

func updateRepoBev(repoBev *repoBeverage, bev model.Beverage) {
	if bev.Type() != "" {
		repoBev.BevType = bev.Type()
	}
	if bev.Brewer() != "" {
		repoBev.Brewer = bev.Brewer()
	}
	if bev.Link() != "" {
		repoBev.Link = bev.Link()
	}
	if bev.Abv() > 0.0 {
		repoBev.Abv = bev.Abv()
	}
	for _, rating := range bev.Ratings() {
		addRepoBevRating(repoBev, rating)
	}
	if repoBev.Attributes == nil {
		repoBev.Attributes = map[string]string{}
	}
	for attr, value := range bev.Attributes() {
		repoBev.Attributes[attr] = value
	}
}

func addRepoBevRating(repoBev *repoBeverage, rating model.Rating) {
	existingRating := findRepoBevRating(repoBev, rating.Source())
	if existingRating != nil {
		existingRating.PercentageRating = rating.PercentageRating()
		return
	}
	repoBev.Ratings = append(repoBev.Ratings, repoRating{
		Source:           rating.Source(),
		PercentageRating: rating.PercentageRating(),
	})
}

func findRepoBevRating(repoBev *repoBeverage, source string) *repoRating {
	for i, _ := range repoBev.Ratings {
		repoRating := &repoBev.Ratings[i]
		if repoRating.Source == source {
			return repoRating
		}
	}
	return nil
}
