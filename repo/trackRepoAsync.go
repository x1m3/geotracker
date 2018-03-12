package repo

import (
	"github.com/x1m3/geotracker/entity"
	"log"
)

// This is a "meta" repo that works using other track repos, using it in an asynchronous way.
// Every time we got a Store() call, it enqueues the operation and returns immediately.
//
// There is a pool of workers waiting for any Request and one of this workers will perform the operation in background.
// The communication with these workers is done with a buffered channel, so, if all workers are busy, the request
// will be enqueued until some worker will go be idle.
//
// The disadvantage of this implementation is that is a "fire and forget". If there is any error storing the request,
// it will fail silently.
type TrackRepoAsync struct {
	repo        Track
	requestChan chan *entity.Track
}

// r Track: Pass the repo to use.
// queueSize int: This is the size of the queue. It limits the amount of request that can be enqueued at any moment,
// reducing pressure over the database.
// nWorkers: Number of workers that will be waiting for request. It increase the pressure over the database.
func NewTrackRepoAsync(r Track, queueSize int, nWorkers int) *TrackRepoAsync {
	repo := &TrackRepoAsync{repo: r}
	repo.requestChan = make(chan *entity.Track, queueSize)
	repo.launchStoreWorkers(nWorkers)
	return repo

}

// Instead of storing the request in the database, it simply enqueues it in a channel and returns ASAP.
func (r *TrackRepoAsync) Store(track *entity.Track) error {
	r.requestChan <- track
	return nil
}

// No workers here. We simply do the request using the underlying repo.
func (r *TrackRepoAsync) GetTracksByDriverAsc(driverID int64) ([]*entity.Track, error) {
	return r.repo.GetTracksByDriverAsc(driverID)
}

// Launches n workers that are blocked reading over the requestChan. Anytime a request arrives from the channel, one
// of the goroutines will get the request and will start processing it. The rest of workers will be idle until they
// gain access to the channel. After performing the store operation, the worker will be blocked again waiting for
// the next request.
func (r *TrackRepoAsync) launchStoreWorkers(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for {
				track := <-r.requestChan
				if err := r.repo.Store(track); err != nil {
					log.Println(err)
				}
			}
		}()
	}
}
