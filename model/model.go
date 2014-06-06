package model

type MenuProvider interface {
	Name() string
	Url() string
	MenuFormat() string
}

type Beverage interface {
	DisplayName() string
	Type() string
	Brewer() string
	SetBrewer(brewer string)
	Link() string
	SetLink(link string)

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

func CreateRating(source string, rating int) Rating {
	return &beverageRating{source: source, rating: rating}
}

// stats

type beverageData struct {
	displayName string
	bevType     string
	brewer      string
	abv         float64
	ratings     []Rating
	link        string
}

func (b *beverageData) HasAbv() bool {
	return b.abv > 0.0
}

func (b *beverageData) Abv() float64 {
	return b.abv
}

func (b *beverageData) SetAbv(abv float64) {
	b.abv = abv
}

func (b *beverageData) HasRating() bool {
	return len(b.ratings) > 0
}

func (b *beverageData) Ratings() []Rating {
	return b.ratings
}

func (b *beverageData) AddRating(rating Rating) {
	b.ratings = append(b.ratings, rating)
}

func (b *beverageData) SetBrewer(brewer string) {
	b.brewer = brewer
}

func (b *beverageData) Brewer() string {
	return b.brewer
}

func (b *beverageData) DisplayName() string {
	return b.displayName
}

func (b *beverageData) Type() string {
	return b.bevType
}

func (b *beverageData) Link() string {
	return b.link
}

func (b *beverageData) SetLink(link string) {
	b.link = link
}

func CreateBeverageBrewer(name string, brewer string) Beverage {
	return &beverageData{displayName: name, brewer: brewer}
}

func CreateBeverage(name string) Beverage {
	return &beverageData{displayName: name}
}

func CreateBeverageAbvTypeRatingLink(name string, abv float64, bevType string,
	rating int, ratingSource string, link string) {
	beverage := CreateBeverage(name)
	beverage.SetAbv(abv)
	beverage.AddRating(CreateRating(ratingSource, rating))
	beverage.SetLink(link)
	return beverage
}

type menuProvider struct {
	name       string
	url        string
	menuFormat string
}

func CreateMenuProvider(name string, url string, format string) MenuProvider {
	return &menuProvider{name: name, url: url, menuFormat: format}
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
