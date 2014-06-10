package syncschedule

import (
	bevsync "github.com/bevly/bevly/sync"
	"github.com/robfig/cron"
	"log"
	"sync"
)

var syncChron *cron.Cron
var syncMutex *sync.Mutex = &sync.Mutex{}

const syncInterval = "31m"

func ScheduleSyncs() {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	// Trigger an immediate sync, then schedule the recurring sync.
	bevsync.TriggerSync()

	if syncChron == nil {
		syncChron = cron.New()
		log.Printf("Scheduling sync every %s\n", syncInterval)
		syncChron.AddFunc("@every "+syncInterval, bevsync.TriggerSync)
		syncChron.Start()
	}
}
