package repo

import "github.com/x1m3/geotracker/entity"

type Track interface {
	Store(track *entity.Track) error
	GetTracksByDriverAsc(driverID int64) ([]*entity.Track, error)
}
