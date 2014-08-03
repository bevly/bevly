package sync

import (
	"github.com/bevly/bevly/fetch/menu"
	"github.com/bevly/bevly/fetch/metadata"
	"github.com/bevly/bevly/repository"
	"log"
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

func (s *Syncer) TriggerSync() {
	log.Println("Triggering beverage sync")
	select {
	case s.SyncChannel <- struct{}{}:
	default:
	}
}

func Sync(repo repository.Repository) []error {
	errors := []error{}
	for _, provider := range repo.MenuProviders() {
		beverages, err := menu.FetchMenu(provider)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		repo.SetBeverageMenu(provider, beverages)
	}

	for _, beverage := range repo.BeveragesNeedingSync() {
		err := metadata.FetchMetadata(beverage)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		repo.SaveBeverage(beverage)
	}
	return errors
}
