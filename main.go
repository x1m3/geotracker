package main

import (
	"database/sql"
	"github.com/x1m3/geotracker/HTTPServer"
	"github.com/x1m3/geotracker/command"
	"github.com/x1m3/geotracker/repo"
	"log"
	"runtime"
	"sync"
)

// Launches a command in another goroutine and increments the waitgroup counter
// when starting, and decrements it when the process finish.
func launchInBackground(function func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		function()
		wg.Done()
	}()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Open a db connection
	db, err := sql.Open("mysql", "xime:@/geotracker?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Let's construct a new httpServer
	router := HTTPServer.NewRouter()
	protocolAdapter := HTTPServer.NewJSONAdapter()
	httpServer := HTTPServer.New(router, protocolAdapter)
	httpServer.RegisterEndpoint("/ping", command.NewPing(), "GET")
	httpServer.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTrackRepoMYSQL(db)), "POST")

	// Server will run in his own goroutine. We need to wait for it to finish
	wg := &sync.WaitGroup{}

	launchInBackground(httpServer.Run, wg)

	wg.Wait()
}
