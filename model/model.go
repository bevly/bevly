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
	return &beverageRating{source: source, percentageRating: rating}
}

// stats

type BeverageData struct {
	displayName string
	bevType     string
	brewer      string
	abv         float64
	ratings     []Rating
	link        string
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
