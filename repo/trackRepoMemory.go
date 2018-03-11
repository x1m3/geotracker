package repo

import (
	"sync"
	"github.com/x1m3/geotracker/entity"
	"errors"
)

type TrackRepoMemory struct {
	sync.RWMutex
	TracksByDriver map[int64][]*entity.Track
}

func NewTracRepoMemory() *TrackRepoMemory {
	return &TrackRepoMemory{TracksByDriver: make(map[int64][]*entity.Track)}
}

func (r *TrackRepoMemory) Store(track *entity.Track) error {
	r.Lock()
	defer r.Unlock()

	if _, found := r.TracksByDriver[track.DriverID]; !found {
		r.TracksByDriver[track.DriverID] = make([]*entity.Track, 0)
	}
	r.TracksByDriver[track.DriverID] = append(r.TracksByDriver[track.DriverID], track)
	return nil
}

func (r *TrackRepoMemory) GetTracksByDriverAsc(driverID int64) ([]*entity.Track, error) {
	var tracks []*entity.Track
	var found bool

	r.RLock()
	r.RUnlock()

	if tracks, found = r.TracksByDriver[driverID]; !found {
		return nil, errors.New("track not found for driver")
	}
	return tracks, nil
}
