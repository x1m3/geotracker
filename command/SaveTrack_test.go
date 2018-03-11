package command_test

import (
	"testing"
	"github.com/x1m3/geotracker/HTTPServer"
	"bytes"
	"github.com/x1m3/geotracker/command"
	"github.com/x1m3/geotracker/repo"
)

func TestAddNewTrack(t *testing.T) {
	testTracks := []string{
		"{ \"latitude\" : 4.5343, \"longitude\" : 3.34324, \"driver_id\": 1 }",
		"{ \"latitude\" : 4.6424, \"longitude\" : 3.35232, \"driver_id\": 1 }",
		"{ \"latitude\" : 4.7534, \"longitude\" : 3.36534, \"driver_id\": 2 }",
		"{ \"latitude\" : 4.8442, \"longitude\" : 3.37321, \"driver_id\": 2 }",
		"{ \"latitude\" : 4.9224, \"longitude\" : 3.38234, \"driver_id\": 2 }",
	}

	memoryRepo := repo.NewTracRepoMemory()

	// This is the command to run
	saveTrackCommand := command.NewSaveTrack(memoryRepo)

	// Lets simulate the creation of a request from a track
	adapter := HTTPServer.NewJSONAdapter()
	for _, jsonTrack := range testTracks {
		request, err := adapter.Decode(bytes.NewBufferString(jsonTrack))
		if err != nil {
			t.Fatalf("Error decoding track. <%s>", err)
		}
		_, err = saveTrackCommand.Call(request)
		if err != nil {
			t.Error(err)
		}
	}

	// Let's inspect the repo to see if all was stored correcty
	if got, expected := len(memoryRepo.TracksByDriver), 2; got!=expected {
		t.Errorf("Wrong number of drivers. Got <%d>, expecting <%d>", got, expected)
	}

	if got, expected := len(memoryRepo.TracksByDriver[1]), 2; got!=expected {
		t.Errorf("Wrong number of tracks for driver 1. Got <%d>, expected <%d>", got, expected)
	}

	if got, expected := len(memoryRepo.TracksByDriver[2]), 3; got!=expected {
		t.Errorf("Wrong number of tracks for driver 2. Got <%d>, expected <%d>", got, expected)
	}
}
