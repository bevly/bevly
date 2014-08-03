package syncschedule

import (
	"github.com/bevly/bevly/repository"
	bevsync "github.com/bevly/bevly/sync"
	"github.com/robfig/cron"
	"log"
)

const syncInterval = "31m"

type SyncScheduler struct {
	syncChron *cron.Cron
	Sync      bevsync.Syncer
}

func CreateSyncScheduler(repo repository.Repository) *SyncScheduler {
	scheduler := SyncScheduler{
		syncChron: cron.New(),
		Sync:      bevsync.CreateSyncer(repo),
	}
	// Trigger an immediate sync, then schedule the recurring sync.
	scheduler.Sync.TriggerSync()

	log.Printf("Scheduling sync every %s\n", syncInterval)
	scheduler.syncChron.AddFunc("@every "+syncInterval, scheduler.Sync.TriggerSync)
	scheduler.syncChron.Start()

	return &scheduler
}
