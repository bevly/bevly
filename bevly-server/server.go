package main

import (
	"github.com/bevly/bevly/http"
	"github.com/bevly/bevly/repository/mongorepo"
	"github.com/bevly/bevly/syncschedule"
	"math/rand"
	"time"
)

func initRng() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	initRng()

	repo := mongorepo.DefaultRepository()
	syncschedule.CreateSyncScheduler(repo)
	http.BeverageServerBlocking(repo)
}
