package command

import (
	"github.com/mitchellh/mapstructure"
	"github.com/x1m3/geotracker/entity"
	"time"
	"github.com/x1m3/geotracker/repo"
)

type trackDTO struct {
	Lat      float64 `mapstructure:"latitude"`
	Lon      float64 `mapstructure:"longitude"`
	DriverID int64   `mapstructure:"driver_id"`
}

type SaveTrackCommand struct {
	repo repo.Track
}

func NewSaveTrack(r repo.Track) *SaveTrackCommand {
	return &SaveTrackCommand{repo: r}
}

func (r *SaveTrackCommand) Call(req Request) (Response, error) {
	if err := r.guardRequest(req); err != nil {
		return nil, err
	}
	dto := trackDTO{}
	if err := mapstructure.Decode(req, &dto); err != nil {
		return nil, err
	}

	coordinate, err := entity.NewCoordinate(dto.Lat, dto.Lon)
	if err != nil {
		return nil, err
	}
	track := entity.NewTrack(coordinate, dto.DriverID, time.Now())

	return "OK", r.repo.Store(track)
}

func (r *SaveTrackCommand) guardRequest(req Request) error {
	if _, found := req["latitude"]; !found {
		return ErrBadRequest
	}
	if _, found := req["longitude"]; !found {
		return ErrBadRequest
	}
	if _, found := req["driver_id"]; !found {
		return ErrBadRequest
	}
	return nil
}
