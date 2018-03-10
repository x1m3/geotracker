package entity

import (
	"testing"
)

func TestNewCoordinate(t *testing.T) {
	testCoords := [][]float64 { {1.1,2.2}, {90, 180}, {-90,180}, {90, -180}, {-90,-180}}

	for _, testValue := range testCoords {
		coord, err := NewCoordinate(testValue[0], testValue[1])
		if err != nil {
			t.Errorf("Expecting a coordinate. Got an error <%s>", err)
		}
		if got, expected := coord.Lat(), testValue[0]; got!=expected {
			t.Errorf("Bad coordinate. Got <%f>, expecting <%f>", got, expected)
		}
		if got, expected := coord.Lon(), testValue[1]; got!=expected {
			t.Errorf("Bad coordinate. Got <%f>, expecting <%f>", got, expected)
		}
	}
}


func TestNewCoordinateBadValue(t *testing.T) {
	testCoords := [][]float64 { {90.1, 10}, {-90.1, 10}, {10, 180.1}, {10, -180.1}}

	for _, testValue := range testCoords {
		_, err := NewCoordinate(testValue[0], testValue[1])
		if err == nil {
			t.Errorf("Expecting an error. Coordinate [%f,%f] is not valid", testValue[0], testValue[1])
		}

	}
}