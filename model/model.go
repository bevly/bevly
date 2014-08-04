package model

type MenuProvider interface {
	Id() string
	Name() string
	Url() string
	MenuFormat() string
}

type Beverage interface {
	DisplayName() string
	Type() string
	SetType(bevType string)

	Name() string
	SetName(name string)

	Description() string
	SetDescription(desc string)

	Brewer() string
	SetBrewer(brewer string)
	Link() string
	SetLink(link string)
	Attribute(name string) string
	Attributes() map[string]string
	SetAttributes(attr map[string]string)
	SetAttribute(name, value string)

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
	displayName string
	name        string
	description string
	bevType     string
	brewer      string
	abv         float64
	attr        map[string]string
	ratings     []Rating
	link        string
}

func (b *BeverageData) Name() string {
	return b.name
}

func (b *BeverageData) SetName(name string) {
	b.name = name
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

func (m *menuProvider) Id() string {
	return m.id
}

func (m *menuProvider) Name() string {
	return m.name
}

func (m *menuProvider) Url() string {
	return m.url
}

func (m *menuProvider) MenuFormat() string {
	return m.menuFormat
}

func (m *menuProvider) String() string {
	return m.id
}
