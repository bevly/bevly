package main

import (
	"github.com/bevly/bevly/http"
	"github.com/bevly/bevly/syncschedule"
)

func main() {
	syncschedule.ScheduleSyncs()
	http.BeverageServerBlocking()
}
