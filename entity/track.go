package entity

import "time"

type Track struct {
	Point      *Coordinate
	DriverID   int64
	ReceivedOn time.Time
}

func NewTrack(c *Coordinate, driverID int64, t time.Time) *Track {
	return &Track{Point: c, DriverID: driverID, ReceivedOn:t}
}

func (t *Track) Equal(t2 *Track) bool {
	return t.Point.Equal(t2.Point) && t.ReceivedOn.Equal(t2.ReceivedOn)
}
