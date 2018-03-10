package entity

import (
	"testing"
	"time"
)

func TestNewTrack(t *testing.T) {
	c, _ := NewCoordinate(1.1, 2.2)
	track := NewTrack(c, 10, time.Now())
	if track.Point.Lat() != 1.1 || track.Point.Lon() != 2.2 {
		t.Error("Bad Coordinate while creating a track")
	}
	if track.DriverID!=10 {
		t.Error("Bad driverID while creating a track")
	}
}

func TestTrackEquals (t *testing.T) {
	aTime := time.Now()
	c1, _ := NewCoordinate(1.1, 2.2)
	c2, _ := NewCoordinate(1.1, 2.2)
	t1 := NewTrack(c1, 1, aTime)
	t2 := NewTrack(c2, 1 ,aTime)

	if !t1.Equal(t2) {
		t.Error("Tracks should be equal")
	}

	c3, _ := NewCoordinate(10, 100)
	t3 := NewTrack(c3, 1, aTime)
	if t3.Equal(t1) {
		t.Error("Tracks should be different")
	}

	t4 := NewTrack(c2, 1 ,aTime.Add(1*time.Millisecond))
	if t4.Equal(t1) {
		t.Error("Tracks should be different")
	}
}
