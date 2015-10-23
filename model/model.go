package model

import "time"

type MenuProvider interface {
	ID() string
	Name() string
	URL() string
	MenuFormat() string
}

type Beverage interface {
	ID() string
	SetID(id string)

	SearchName() string
	DisplayName() string
	SetDisplayName(name string)
	SetSearchName(name string)
	Type() string
	SetType(bevType string)

	Name() string
	SetName(name string)

	NeedSync() bool
	SetNeedSync(needSync bool)

	Description() string
	SetDescription(desc string)

	AccuracyScore() int
	SetAccuracyScore(accuracy int)

	Brewer() string
	SetBrewer(brewer string)
	Link() string
	SetLink(link string)
	Attribute(name string) string
	Attributes() map[string]string
	SetAttributes(attr map[string]string)
	SetAttribute(name, value string)
	SyncTime() time.Time
	SetSyncTime(newtime time.Time)

	BeverageStats
}

type BeverageStats interface {
	HasAbv() bool
	Abv() float64
	SetAbv(abv float64)

	HasRating() bool
	Ratings() []Rating
	AddRating(rating Rating)
	ClearRatings()
}

type Rating interface {
	Source() string
	PercentageRating() int
	SetPercentageRating(rating int)
}

// rating

type beverageRating struct {
	source           string
	percentageRating int
}

func (r *beverageRating) Source() string {
	return r.source
}

func (r *beverageRating) PercentageRating() int {
	return r.percentageRating
}

func (r *beverageRating) SetPercentageRating(rating int) {
	if rating > 0 {
		r.percentageRating = rating
	}
}

func CreateRating(source string, rating int) Rating {
	return &beverageRating{source: source, percentageRating: rating}
}

// stats

type BeverageData struct {
	id            string
	accuracyScore int
	displayName   string
	name          string
	description   string
	bevType       string
	brewer        string
	abv           float64
	attr          map[string]string
	ratings       []Rating
	link          string
	syncTime      time.Time
	needSync      bool
}

func (b *BeverageData) ID() string {
	return b.id
}

func (b *BeverageData) SetID(id string) {
	b.id = id
}

func (b *BeverageData) Name() string {
	return b.name
}

func (b *BeverageData) SetName(name string) {
	b.name = name
}

func (b *BeverageData) SearchName() string {
	searchName := b.Attribute("SearchName")
	if searchName != "" {
		return searchName
	}
	return b.DisplayName()
}

func (b *BeverageData) SetSearchName(name string) {
	b.SetAttribute("SearchName", name)
}

func (b *BeverageData) AccuracyScore() int {
	return b.accuracyScore
}

func (b *BeverageData) SetAccuracyScore(accuracy int) {
	b.accuracyScore = accuracy
}

func (b *BeverageData) NeedSync() bool {
	return b.needSync
}

func (b *BeverageData) SetNeedSync(needSync bool) {
	b.needSync = needSync
}

func (b *BeverageData) SyncTime() time.Time {
	return b.syncTime
}

func (b *BeverageData) SetSyncTime(newTime time.Time) {
	b.syncTime = newTime
}

func (b *BeverageData) Description() string {
	return b.description
}

func (b *BeverageData) SetDescription(description string) {
	b.description = description
}

func (b *BeverageData) Attribute(name string) string {
	if b.attr == nil {
		return ""
	}
	return b.attr[name]
}

func (b *BeverageData) Attributes() map[string]string {
	return b.attr
}

func (b *BeverageData) SetAttributes(attr map[string]string) {
	b.attr = attr
}

func (b *BeverageData) SetAttribute(name, value string) {
	if b.attr == nil {
		b.attr = map[string]string{}
	}
	b.attr[name] = value
}

func (b *BeverageData) HasAbv() bool {
	return b.abv > 0.0
}

func (b *BeverageData) Abv() float64 {
	return b.abv
}

func (b *BeverageData) SetAbv(abv float64) {
	b.abv = abv
}

func (b *BeverageData) HasRating() bool {
	return len(b.ratings) > 0
}

func (b *BeverageData) Ratings() []Rating {
	return b.ratings
}

func (b *BeverageData) AddRating(rating Rating) {
	if rating.PercentageRating() == 0 {
		return
	}
	for _, existingRating := range b.ratings {
		if existingRating.Source() == rating.Source() {
			existingRating.SetPercentageRating(rating.PercentageRating())
			return
		}
	}
	b.ratings = append(b.ratings, rating)
}

func (b *BeverageData) ClearRatings() {
	b.ratings = []Rating{}
}

func (b *BeverageData) SetBrewer(brewer string) {
	b.brewer = brewer
}

func (b *BeverageData) Brewer() string {
	return b.brewer
}

func (b *BeverageData) DisplayName() string {
	return b.displayName
}

func (b *BeverageData) SetDisplayName(name string) {
	b.displayName = name
}

func (b *BeverageData) Type() string {
	return b.bevType
}

func (b *BeverageData) SetType(bevType string) {
	b.bevType = bevType
}

func (b *BeverageData) Link() string {
	return b.link
}

func (b *BeverageData) SetLink(link string) {
	b.link = link
}

func (b *BeverageData) String() string {
	return b.DisplayName()
}

func CreateBeverageBrewer(name string, brewer string) Beverage {
	return &BeverageData{displayName: name, brewer: brewer}
}

func CreateBeverage(name string) Beverage {
	return &BeverageData{displayName: name}
}

func CreateBeverageAbvTypeRatingLink(name string, abv float64, bevType string,
	rating int, ratingSource string, link string) Beverage {
	beverage := CreateBeverage(name)
	beverage.SetType(bevType)
	beverage.SetAbv(abv)
	beverage.AddRating(CreateRating(ratingSource, rating))
	beverage.SetLink(link)
	return beverage
}

type menuProvider struct {
	id         string
	name       string
	url        string
	menuFormat string
}

func CreateMenuProvider(id string, name string, url string, format string) MenuProvider {
	return &menuProvider{id: id, name: name, url: url, menuFormat: format}
}

func (m *menuProvider) ID() string {
	return m.id
}

func (m *menuProvider) Name() string {
	return m.name
}

func (m *menuProvider) URL() string {
	return m.url
}

func (m *menuProvider) MenuFormat() string {
	return m.menuFormat
}

func (m *menuProvider) String() string {
	return m.id
}
