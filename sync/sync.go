package sync

import (
	"github.com/bevly/bevly/fetch/menu"
	"github.com/bevly/bevly/fetch/metadata"
	"github.com/bevly/bevly/model"
	"github.com/bevly/bevly/repository"
	"log"
	"math/rand"
	"time"
)

type Syncer struct {
	Repo        repository.Repository
	SyncChannel chan struct{}
}

func CreateSyncer(repo repository.Repository) Syncer {
	syncer := Syncer{
		Repo:        repo,
		SyncChannel: make(chan struct{}),
	}
	syncer.startSyncJob()
	return syncer
}

func (s *Syncer) startSyncJob() {
	go s.syncJob()
}

func (s *Syncer) syncJob() {
	for {
		<-s.SyncChannel
		Sync(s.Repo)
	}
}

func (s *Syncer) TriggerSync(blocking bool) {
	log.Println("Triggering beverage sync")
	if blocking {
		s.SyncChannel <- struct{}{}
	} else {
		select {
		case s.SyncChannel <- struct{}{}:
		default:
		}
	}
}

func Sync(repo repository.Repository) []error {
	errors := []error{}
	log.Println("Syncing all providers")
	for _, provider := range repo.MenuProviders() {
		log.Printf("Syncing provider: %s\n", provider)
		beverages, err := menu.FetchMenu(provider)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		priorBeverages := repo.ProviderBeverages(provider)
		SetBeverageDiscoverTimes(provider, beverages, priorBeverages)
		repo.SetBeverageMenu(provider, beverages)
	}

	for _, beverage := range repo.BeveragesNeedingSync() {
		dur := randSleepInterval()
		log.Printf("Sleeping %dms before metadata fetch for %s\n",
			int64(dur)/1000, beverage)
		time.Sleep(dur)
		beverage.SetNeedSync(false)
		err := metadata.FetchMetadata(beverage)
		if err != nil {
			errors = append(errors, err)
		}
		if beverage.NeedSync() {
			repo.SaveBeverage(beverage)
		}
	}
	return errors
}

// SetBeverageDiscoverTimes sets the discovery time for each beverage
// that was not in provider's prior menu.
func SetBeverageDiscoverTimes(provider model.MenuProvider, bevs []model.Beverage, priorBevs []model.Beverage) {
	seenBevs := beverageNameMap(priorBevs)
	syncTimeISO := time.Now().Format(time.RFC3339)
	for _, bev := range bevs {
		if !seenBevs[bev.DisplayName()] {
			// This beverage was just added.
			bev.SetAttribute(provider.Id()+"MenuAt", syncTimeISO)
		}
	}
}

func beverageNameMap(bevs []model.Beverage) map[string]bool {
	result := map[string]bool{}
	for _, bev := range bevs {
		result[bev.DisplayName()] = true
	}
	return result
}

func randSleepInterval() time.Duration {
	return time.Duration(randRange(3500, 18500)) * time.Millisecond
}

func randRange(low, hi int) int {
	return low + rand.Intn(hi-low+1)
}
