package mongorepo

import (
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/policy"
	"github.com/bevly/bevly/repository"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"sync"
	"time"
)

const BeverageDbName = "bevly"
const BeverageFetchLimit = 1000

type mongoRepo struct {
	session     *mgo.Session
	db          *mgo.Database
	providers   *mgo.Collection
	beverages   *mgo.Collection
	initialized bool
	mutex       sync.Mutex
}

type repoProvider struct {
	Id          bson.ObjectId   `bson:"_id"`
	ProviderId  string          `bson:"providerId"`
	Name        string          `bson:"name"`
	Url         string          `bson:"url"`
	MenuFormat  string          `bson:"menuFormat"`
	BeverageIds []bson.ObjectId `bson:"beverageIds"`
}

type repoBeverage struct {
	Id          bson.ObjectId `bson:"_id"`
	DisplayName string        `bson:"displayName"`
	BevType     string        `bson:"bevType"`
	Brewer      string        `bson:"brewer"`
	Abv         float64       `bson:"abv"`
	Ratings     []repoRating  `bson:"ratings"`
	Link        string        `bson:"link"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
}

type repoRating struct {
	Source           string `bson:"source"`
	PercentageRating int    `bson:"percentageRating"`
}

var globalSession = mongoRepo{}

func DefaultRepository() repository.Repository {
	globalSession.MustInit()
	return &globalSession
}

func mongoServer() string {
	envMongoServer := os.Getenv("MONGO_HOST")
	if envMongoServer != "" {
		return envMongoServer
	}
	return "localhost"
}

func (repo *mongoRepo) MustInit() {
	err := repo.Init()
	if err != nil {
		panic(err)
	}
}

func (repo *mongoRepo) Init() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if repo.initialized {
		return nil
	}

	var err error
	repo.session, err = mgo.Dial(mongoServer())
	if err != nil {
		return err
	}

	repo.db = repo.session.DB(BeverageDbName)
	repo.providers = repo.db.C("providers")
	repo.beverages = repo.db.C("beverages")
	repo.initialized = true
	return nil
}

func (repo *mongoRepo) Purge() {
	err := repo.providers.DropCollection()
	if err != nil {
		log.Printf("Error purging providers collection: %s", err)
	}
	err = repo.beverages.DropCollection()
	if err != nil {
		log.Printf("Error purging beverages collection: %s", err)
	}
}

func (repo *mongoRepo) GarbageCollect() {
	referencedBeverageIds := repo.beverageIdsReferencedInMenus()
	discardThresholdTime := policy.BeverageDiscardThresholdTime()
	changes, err := repo.beverages.RemoveAll(
		bson.M{
			"_id":       bson.M{"$nin": referencedBeverageIds},
			"updatedAt": bson.M{"$lt": discardThresholdTime},
		})
	if err != nil {
		log.Printf("GarbageCollect(older:%v): error %s\n", discardThresholdTime, err)
	} else {
		log.Printf("GarbageCollect(older:%v): removed %d beverages\n", discardThresholdTime, changes.Removed)
	}
}

func (repo *mongoRepo) MenuProviders() []model.MenuProvider {
	return repository.StubRepository().MenuProviders()
}

func (repo *mongoRepo) ProviderById(id string) model.MenuProvider {
	return repository.StubRepository().ProviderById(id)
}

func (repo *mongoRepo) ProviderBeverages(prov model.MenuProvider) []model.Beverage {
	provider, err := repo.findProvider(prov)
	if err != nil {
		return []model.Beverage{}
	}
	result, err := repo.lookupBeveragesByIds(provider.BeverageIds)
	if err != nil {
		log.Printf("Could not look up beverages for provider %s (%s) with ids: %v\n",
			prov.Name(), prov.Id(), provider.BeverageIds)
	}
	return result
}

func (repo *mongoRepo) ProviderIdBeverages(id string) []model.Beverage {
	return repo.ProviderBeverages(repo.ProviderById(id))
}

func (repo *mongoRepo) BeveragesNeedingSync() []model.Beverage {
	referencedBeverageIds := repo.beverageIdsReferencedInMenus()
	staleUpdateTime := policy.BeverageResyncThresholdTime()
	var beverages []repoBeverage
	err := repo.beverages.Find(
		bson.M{
			"_id": bson.M{"$in": referencedBeverageIds},
			"$or": []interface{}{
				bson.M{"syncTime": nil},
				bson.M{"syncTime": bson.M{"$lt": staleUpdateTime}},
			},
		}).All(&beverages)
	if err != nil {
		log.Printf("Error looking up beverages needing sync: %s\n", err)
	}
	log.Printf("Found %d beverages needing sync\n", len(beverages))
	return repoBeverageModels(beverages)
}

func (repo *mongoRepo) SetBeverageMenu(prov model.MenuProvider, beverages []model.Beverage) {
	beverageIds, err := repo.saveBeverages(beverages)
	if err != nil {
		log.Printf("Failed to save beverages for %s: %v", prov.Name(), err)
		return
	}
	err = repo.saveProviderMenu(prov, beverageIds)
	if err != nil {
		log.Printf("Failed to save provider menu for %s: %s", prov.Name(), err)
	}
}

func (repo *mongoRepo) SaveBeverage(beverage model.Beverage) {
	_, err := repo.saveBeverage(beverage)
	if err != nil {
		log.Printf("SaveBeverage(%s) failed: %s", beverage, err)
	}
}

func (repo *mongoRepo) BeverageByName(name string) model.Beverage {
	repoBev, err := repo.findBeverageByName(name)
	if err != nil {
		return nil
	}
	return repoBeverageModel(repoBev)
}

func (repo *mongoRepo) saveProviderMenu(prov model.MenuProvider, beverageIds []bson.ObjectId) error {
	provider, err := repo.findProvider(prov)
	if err == nil { // menu exists
		provider.BeverageIds = beverageIds
		_, err = repo.providers.UpsertId(provider.Id, provider)
		return err
	} else {
		provider = &repoProvider{
			Id:          bson.NewObjectId(),
			ProviderId:  prov.Id(),
			Name:        prov.Name(),
			Url:         prov.Url(),
			MenuFormat:  prov.MenuFormat(),
			BeverageIds: beverageIds,
		}
		return repo.providers.Insert(provider)
	}
}

func (repo *mongoRepo) saveBeverages(beverages []model.Beverage) ([]bson.ObjectId, error) {
	beverageIds := make([]bson.ObjectId, 0, len(beverages))

	errors := &compositeError{}
	for _, beverage := range beverages {
		id, err := repo.saveBeverage(beverage)
		if err != nil {
			errors.Add(err)
			continue
		}
		beverageIds = append(beverageIds, id)
	}
	if errors.IsError() {
		return beverageIds, errors
	}
	return beverageIds, nil
}

func (repo *mongoRepo) saveBeverage(beverage model.Beverage) (bson.ObjectId, error) {
	// Update or insert
	repoBev, err := repo.findBeverageByName(beverage.DisplayName())

	updateTime := time.Now()
	if err == nil { // found existing object
		updateRepoBev(repoBev, beverage)
		log.Printf("Updating beverage %s with id %s",
			repoBev.DisplayName, repoBev.Id)
		repoBev.UpdatedAt = updateTime
		_, err := repo.beverages.UpsertId(repoBev.Id, repoBev)
		if err != nil {
			return bson.ObjectId(""), err
		}
		return repoBev.Id, nil
	} else {
		repoBev = beverageModelToRepo(beverage)
		repoBev.Id = bson.NewObjectId()
		repoBev.UpdatedAt = updateTime
		log.Printf("Inserting beverage %s with id %s",
			repoBev.DisplayName, repoBev.Id)
		err = repo.beverages.Insert(repoBev)
		if err != nil {
			return bson.ObjectId(""), err
		}
		return repoBev.Id, nil
	}
}

func (repo *mongoRepo) findBeverageByName(name string) (*repoBeverage, error) {
	repoBev := &repoBeverage{}
	err := repo.beverages.Find(bson.M{"displayName": name}).One(repoBev)
	if err != nil {
		return nil, err
	}
	return repoBev, nil
}

func (repo *mongoRepo) lookupBeveragesByIds(ids []bson.ObjectId) ([]model.Beverage, error) {
	var beverages []repoBeverage
	err := repo.beverages.Find(bson.M{"_id": bson.M{"$in": ids}}).Limit(BeverageFetchLimit).All(&beverages)
	if err != nil {
		return nil, err
	}
	return repoBeverageModels(beverages), nil
}

func (repo *mongoRepo) findProvider(prov model.MenuProvider) (*repoProvider, error) {
	provider := repoProvider{}
	err := repo.providers.Find(providerQuery(prov)).One(&provider)
	return &provider, err
}

func providerQuery(prov model.MenuProvider) bson.M {
	return bson.M{"providerId": prov.Id()}
}

func (repo *mongoRepo) beverageIdsReferencedInMenus() []bson.ObjectId {
	providerIter := repo.providers.Find(nil).Iter()
	provider := repoProvider{}
	referencedBeverageIds := map[bson.ObjectId]bool{}
	for providerIter.Next(&provider) {
		for _, beverageId := range provider.BeverageIds {
			referencedBeverageIds[beverageId] = true
		}
	}

	result := make([]bson.ObjectId, 0, len(referencedBeverageIds))
	for id, _ := range referencedBeverageIds {
		result = append(result, id)
	}
	return result
}
