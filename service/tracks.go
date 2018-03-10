package service

import (
	"github.com/x1m3/geotracker/repo"
	"github.com/x1m3/geotracker/entity"
)

func StoreATrack(repo repo.Track, t *entity.Track) error {
	return repo.Store(t)
}

func GetTracksByDriverASC (repo repo.Track, driverID int64) ([]*entity.Track, error) {
	return repo.GetTracksByDriverAsc(driverID)
}
