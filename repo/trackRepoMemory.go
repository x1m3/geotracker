package repo

import (
	"sync"
	"github.com/x1m3/geotracker/entity"
	"errors"
)

type TrackRepoMemory struct {
	sync.RWMutex
	tracksByDriver map[int64][]*entity.Track
}

func NewTracRepoMemory() *TrackRepoMemory {
	return &TrackRepoMemory{tracksByDriver: make(map[int64][]*entity.Track)}
}

func (r *TrackRepoMemory) Store(track *entity.Track) error {
	r.Lock()
	defer r.Unlock()

	if _, found := r.tracksByDriver[track.DriverID]; !found {
		r.tracksByDriver[track.DriverID] = make([]*entity.Track, 0)
	}
	r.tracksByDriver[track.DriverID] = append(r.tracksByDriver[track.DriverID], track)
	return nil
}

func (r *TrackRepoMemory) GetTracksByDriverAsc(driverID int64) ([]*entity.Track, error) {
	var tracks []*entity.Track
	var found bool

	r.RLock()
	r.RUnlock()

	if tracks, found = r.tracksByDriver[driverID]; !found {
		return nil, errors.New("track not found for driver")
	}
	return tracks, nil
}
