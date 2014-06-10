package sync

import (
	"github.com/bevly/bevly/fetch/menu"
	"github.com/bevly/bevly/fetch/metadata"
	"github.com/bevly/bevly/repository"
	"log"
)

var syncChannel chan bool = make(chan bool)

func Sync() []error {
	errors := []error{}
	repo := repository.DefaultRepository()
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

func init() {
	go syncJob()
}

func TriggerSync() {
	log.Println("Triggering beverage sync")
	select {
	case syncChannel <- true:
	default:
	}
}

func syncJob() {
	for {
		<-syncChannel
		Sync()
	}
}
