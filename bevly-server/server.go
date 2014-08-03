package main

import (
	"github.com/bevly/bevly/http"
	"github.com/bevly/bevly/repository/mongorepo"
	"github.com/bevly/bevly/syncschedule"
)

func main() {
	repo := mongorepo.DefaultRepository()
	syncschedule.CreateSyncScheduler(repo)
	http.BeverageServerBlocking(repo)
}
