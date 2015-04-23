package main

import (
	"github.com/bevly/bevly/http"
	"github.com/bevly/bevly/repository/mongorepo"
	"github.com/bevly/bevly/syncschedule"
	"log"
	"math/rand"
	"os"
	"time"
)

func initRng() {
	rand.Seed(time.Now().UnixNano())
}

func syncEnabled() bool {
	return os.Getenv("BEVLY_SYNC_DISABLE") == ""
}

func main() {
	initRng()

	repo := mongorepo.DefaultRepository()

	if syncEnabled() {
		log.Println("Creating sync scheduler")
		syncschedule.CreateSyncScheduler(repo)
	} else {
		log.Println("Sync is disabled")
	}

	http.BeverageServerBlocking(repo)
}
