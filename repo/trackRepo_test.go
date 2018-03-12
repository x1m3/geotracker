package repo

import (
	"github.com/x1m3/geotracker/entity"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

func initTestData() map[int64]entity.TrackList {
	var driver int64
	const N_DRIVERS = 10
	const TRACKS_PER_DRIVER = 500

	randomTracks := make(map[int64]entity.TrackList, 0)

	for driver = 0; driver < N_DRIVERS; driver++ {
		c, _ := entity.NewCoordinate(0.0, 0.0)
		newLat, newLon := 0.0, 0.0
		t := time.Now()
		for i := 0; i < TRACKS_PER_DRIVER; i++ {
			c, _ = entity.NewCoordinate(newLat, newLon)
			randomTracks[driver] = append(randomTracks[driver], entity.NewTrack(c, driver, t))
			newLat = c.Lat() + float64(rand.Intn(10)-5)/10000
			newLon = c.Lon() + float64(rand.Intn(10)-5)/10000
			t = t.Add(1 * time.Second)
		}
	}
	return randomTracks
}

// A helper function that insert all tracks for a driver
func insertDriverTracks(repo Track, tracks entity.TrackList, t *testing.T) {
	for _, track := range tracks {
		if err := repo.Store(track); err != nil {
			t.Errorf("Error inserting track. <%s>", err)
		}
	}
}

func TestTrackRepoSuite(t *testing.T) {
	fp1, err := ioutil.TempFile(os.TempDir(), "trakRepo")
	if err != nil {
		t.Fatalf("Cannot create temp file for test <%s>", err)
	}
	fp2, err := ioutil.TempFile(os.TempDir(), "trakRepo")
	if err != nil {
		t.Fatalf("Cannot create temp file for test <%s>", err)
	}
	defer os.Remove(fp1.Name())
	defer os.Remove(fp2.Name())

	repos := make([]Track, 0)
	repos = append(repos, NewTracRepoFile(fp1))
	repos = append(repos, NewTracRepoMemory())
	repos = append(repos, NewTrackRepoAsync(NewTracRepoFile(fp2), 1000, 100))
	repos = append(repos, NewTrackRepoAsync(NewTracRepoMemory(), 1000, 100))

	/*
		db, err := sql.Open("mysql", "xime:@/geotracker?parseTime=true")
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		repos = append(repos, NewTrackRepoAsync(NewTrackRepoMYSQL(db), 1000, 10))
	*/

	randomTracks := initTestData()

	for _, r := range repos { // Test all available repos
		testTrackRepo(r, randomTracks, t)
	}

}

func testTrackRepo(repo Track, randomTracks map[int64]entity.TrackList, t *testing.T) {

	// Let's insert all tracks in repo, all drivers at the "same" time, one goroutine per driver
	wg := sync.WaitGroup{}
	for _, driverTracks := range randomTracks {
		wg.Add(1)
		go func(driverTracks entity.TrackList) {
			insertDriverTracks(repo, driverTracks, t)
			wg.Done()
		}(driverTracks)
	}
	wg.Wait()

	time.Sleep(1 * time.Second)

	// Let's read all driver tracks from repo, to confirm that all was stored
	for driver, tracks := range randomTracks {
		repoTracks, err := repo.GetTracksByDriverAsc(driver)
		if err != nil {
			t.Fatalf("Error reading tracks. <%s>", err)
		}
		if got, expected := len(repoTracks), len(tracks); got != expected {
			t.Errorf("Error reading tracks for driver <%d>. Len differs. got <%d>, expecting <%d>", driver, got, expected)
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
