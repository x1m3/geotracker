package main

import (
	"runtime"
	"github.com/x1m3/geotracker/HTTPServer"
	"sync"
	"github.com/x1m3/geotracker/repo"
	"github.com/x1m3/geotracker/command"
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

	// Let's construct a new httpServer
	router := HTTPServer.NewRouter()
	protocolAdapter := HTTPServer.NewJSONAdapter()
	httpServer := HTTPServer.New(router, protocolAdapter)
	httpServer.RegisterEndpoint("/ping", command.NewPing(), "GET")
	httpServer.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTracRepoMemory()), "POST")

	// Server will run in his own goroutine. We need to wait for it to finish
	wg := &sync.WaitGroup{}

	launchInBackground(httpServer.Run, wg)

	wg.Wait()
}



