package repo

import (
	"errors"
	"github.com/x1m3/geotracker/entity"
	"sort"
	"sync"
)

type TrackRepoMemory struct {
	sync.RWMutex
	TracksByDriver map[int64]entity.TrackList
}

func NewTracRepoMemory() *TrackRepoMemory {
	return &TrackRepoMemory{TracksByDriver: make(map[int64]entity.TrackList)}
}

func (r *TrackRepoMemory) Store(track *entity.Track) error {
	r.Lock()
	defer r.Unlock()

	if _, found := r.TracksByDriver[track.DriverID]; !found {
		r.TracksByDriver[track.DriverID] = make(entity.TrackList, 0)
	}
	r.TracksByDriver[track.DriverID] = append(r.TracksByDriver[track.DriverID], track)
	return nil
}

func (r *TrackRepoMemory) GetTracksByDriverAsc(driverID int64) (entity.TrackList, error) {
	var tracks entity.TrackList
	var found bool

	r.RLock()
	defer r.RUnlock()

	if tracks, found = r.TracksByDriver[driverID]; !found {
		return nil, errors.New("track not found for driver")
	}
	sort.Sort(entity.ByDate(tracks))
	return tracks, nil
}
