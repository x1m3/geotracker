package repo

import (
	"testing"
	"io/ioutil"
	"os"
	"github.com/x1m3/geotracker/entity"
	"math/rand"
	"sync"
	"time"
)

func initTestData() map[int64][]*entity.Track {
	var driver int64
	const N_DRIVERS = 10
	const TRACKS_PER_DRIVER = 1000

	randomTracks := make(map[int64][]*entity.Track, 0)

	for driver = 0; driver < N_DRIVERS; driver++ {
		c, _ := entity.NewCoordinate(0.0, 0.0)
		for i := 0; i < TRACKS_PER_DRIVER; i++ {
			randomTracks[driver] = append(randomTracks[driver], entity.NewTrack(c, driver, time.Now()))
			newLat := c.Lat() + float64(rand.Intn(10)-5)/10000
			newLon := c.Lon() + float64(rand.Intn(10)-5)/10000
			c, _ = entity.NewCoordinate(newLat, newLon)
		}
	}
	return randomTracks
}

// A helper function that insert all tracks for a driver
func insertDriverTracks(repo TrackRepo, tracks []*entity.Track, t *testing.T) {
	for _, track := range tracks {
		if err := repo.Store(track); err != nil {
			t.Errorf("Error inserting track. <%s>", err)
		}
	}
}

func TestTrackRepoSuite(t *testing.T) {
	fp, err := ioutil.TempFile(os.TempDir(), "trakRepo")
	if err != nil {
		t.Fatalf("Cannot create temp file for test <%s>", err)
	}
	defer os.Remove(fp.Name())
	r := NewTracRepoFile(fp)

	randomTracks := initTestData()

	testTrackRepo(r, randomTracks, t)

}

func testTrackRepo(repo TrackRepo, randomTracks map[int64][]*entity.Track, t *testing.T) {

	// Let's insert all tracks in repo, all drivers at the "same" time, one goroutine per driver
	wg := sync.WaitGroup{}
	for _, driverTracks := range randomTracks {
		wg.Add(1)
		go func(driverTracks []*entity.Track) {
			insertDriverTracks(repo, driverTracks, t)
			wg.Done()
		}(driverTracks)
	}
	wg.Wait()

	// Let's read all driver tracks from repo, to confirm that all was stored
	for driver, tracks := range randomTracks {
		repoTracks, err := repo.GetTracksByDriver(driver)
		if err != nil {
			t.Fatalf("Error reading tracks. <%s>", err)
		}
		if got, expected := len(repoTracks), len(tracks); got != expected {
			t.Fatalf("Error reading tracks for driver <%d>. Len differs. got <%d>, expecting <%d>", driver, got, expected)
		}

		for i := 0; i < len(tracks); i++ {
			t1 := tracks[i]
			t2 := repoTracks[i]
			if !t1.Equal(t2) {
				t.Errorf("Entities differ. got [lat:%v, lon:%v, receivedOn:%v], expected [lat:%v, lon:%v, receivedOn:%v]",
					t1.Point.Lat(),
					t1.Point.Lon(),
					t1.ReceivedOn,
					t2.Point.Lat(),
					t2.Point.Lon(),
					t2.ReceivedOn)
			}
		}
	}

}
