package sync

import (
	"github.com/bevly/bevly/fetch/menu"
	"github.com/bevly/bevly/fetch/metadata"
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
		repo.SetBeverageMenu(provider, beverages)
	}

	for _, beverage := range repo.BeveragesNeedingSync() {
		dur := randSleepInterval()
		log.Printf("Sleeping %dms before metadata fetch for %s\n",
			int64(dur)/1000, beverage)
		time.Sleep(dur)
		err := metadata.FetchMetadata(beverage)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		repo.SaveBeverage(beverage)
	}
	return errors
}

func randSleepInterval() time.Duration {
	return time.Duration(randRange(1500, 12500)) * time.Millisecond
}

func randRange(low, hi int) int {
	return low + rand.Intn(hi-low+1)
}
