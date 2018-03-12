package entity

import (
	"math"
	"time"
)

type TrackList []*Track

type ByDate TrackList

func (list ByDate) Len() int           { return len(list) }
func (list ByDate) Less(i, j int) bool { return list[i].ReceivedOn.Before(list[j].ReceivedOn) }
func (list ByDate) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

type Track struct {
	Point      *Coordinate
	DriverID   int64
	ReceivedOn time.Time
}

func NewTrack(c *Coordinate, driverID int64, t time.Time) *Track {
	return &Track{Point: c, DriverID: driverID, ReceivedOn: t}
}

func (t *Track) Equal(t2 *Track) bool {
	return t.Point.Equal(t2.Point) &&
		time.Duration(math.Abs(float64(t.ReceivedOn.Sub(t2.ReceivedOn)))) <= 10*time.Microsecond
}
