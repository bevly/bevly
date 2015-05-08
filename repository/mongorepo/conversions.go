package mongorepo

import (
	"github.com/bevly/bevly/model"

	"encoding/hex"
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
	bev.SetId(hex.EncodeToString([]byte(repoBev.Id)))
	bev.SetType(repoBev.BevType)
	bev.SetName(repoBev.Name)
	bev.SetDescription(repoBev.Description)
	bev.SetBrewer(repoBev.Brewer)
	bev.SetLink(repoBev.Link)
	bev.SetAbv(repoBev.Abv)
	bev.SetSyncTime(repoBev.SyncTime)
	bev.SetAttributes(repoBev.Attributes)
	bev.SetAccuracyScore(repoBev.AccuracyScore)
	for _, rating := range repoBev.Ratings {
		bev.AddRating(model.CreateRating(rating.Source, rating.PercentageRating))
	}
	return bev
}

func beverageModelToRepo(bev model.Beverage) *repoBeverage {
	repoBev := &repoBeverage{}
	repoBev.DisplayName = bev.DisplayName()
	repoBev.Name = bev.Name()
	repoBev.Description = bev.Description()
	repoBev.BevType = bev.Type()
	repoBev.Brewer = bev.Brewer()
	repoBev.Link = bev.Link()
	repoBev.Abv = bev.Abv()
	repoBev.Attributes = bev.Attributes()
	repoBev.SyncTime = bev.SyncTime()
	repoBev.AccuracyScore = bev.AccuracyScore()

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
	overwrite := bev.AccuracyScore() >= repoBev.AccuracyScore

	set := func(oldVal *string, newVal string) {
		if newVal != "" && (overwrite || *oldVal == "") {
			*oldVal = newVal
		}
	}

	if !bev.SyncTime().IsZero() {
		repoBev.SyncTime = bev.SyncTime()
	}
	set(&repoBev.BevType, bev.Type())
	set(&repoBev.Name, bev.Name())
	set(&repoBev.Description, bev.Description())
	set(&repoBev.Brewer, bev.Brewer())
	set(&repoBev.Link, bev.Link())
	if bev.Abv() > 0.0 && (overwrite || repoBev.Abv == 0.0) {
		repoBev.Abv = bev.Abv()
	}
	for _, rating := range bev.Ratings() {
		addRepoBevRating(repoBev, rating)
	}
	if repoBev.Attributes == nil {
		repoBev.Attributes = map[string]string{}
	}
	for attr, value := range bev.Attributes() {
		if value != "" {
			repoBev.Attributes[attr] = value
		}
	}
	if bev.AccuracyScore() > repoBev.AccuracyScore {
		repoBev.AccuracyScore = bev.AccuracyScore()
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
