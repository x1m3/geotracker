package entity

import (
	"math"
	"errors"
)

type Coordinate struct {
	lat float64
	lon float64
}

func NewCoordinate(lat float64, lon float64) (*Coordinate, error) {
	if math.Abs(lat)>90 {
		return nil, errors.New("wrong latitude")
	}
	if math.Abs(lon)>180 {
		return nil, errors.New("wrong longitude")
	}
	return &Coordinate{lat:lat, lon:lon},nil
}

func (c *Coordinate) Lat() float64 {
	return c.lat
}


func (c *Coordinate) Lon() float64 {
	return c.lon
}

func (c *Coordinate) Equal(other *Coordinate) bool {
	return c.lat== other.lat && c.lon == other.lon
}
