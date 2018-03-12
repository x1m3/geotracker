package HTTPServer

import (
	"bytes"
	"encoding/json"
	"github.com/x1m3/geotracker/command"
	"github.com/x1m3/geotracker/repo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPServer_Ping(t *testing.T) {
	server := New(NewRouter(), NewJSONAdapter(), "", 80)
	server.RegisterEndpoint("/ping", command.NewPing(), "GET")

	testServer := httptest.NewServer(server.httpServer.Handler)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/ping")
	if err != nil {
		t.Errorf("Error pinging HTTPServer. <%s>", err)
	}

	if got, expected := resp.StatusCode, http.StatusOK; got != expected {
		t.Errorf("Bad Status Code. Got <%v>, expecting <%v>", got, expected)
	}

	if got, expected := resp.Header.Get("Content-Type"), "application/json"; got != expected {
		t.Errorf("Bad Content Type. Got <%s>, expecting <%s>", got, expected)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading body <%s>", err)
	}

	got := ""
	if err := json.Unmarshal(body, &got); err != nil {
		t.Errorf("Error decoding json response <%s>", err)
	}

	if expected := "pong"; got != expected {
		t.Errorf("Wrong body response for /ping. Got <%s>, expecting <%s>", got, expected)
	}
}

func TestHTTPServerNotFound(t *testing.T) {
	server := New(NewRouter(), NewJSONAdapter(), "", 80)

	testServer := httptest.NewServer(server.httpServer.Handler)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/this/url/doesnt/exists")
	if err != nil {
		t.Errorf("Error pinging HTTPServer. <%s>", err)
	}

	if got, expected := resp.StatusCode, http.StatusNotFound; got != expected {
		t.Errorf("Bad Status Code. Got <%v>, expecting <%v>", got, expected)
	}

	if got, expected := resp.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != expected {
		t.Errorf("Bad Content Type. Got <%s>, expecting <%s>", got, expected)
	}
}

func TestStoreAtrack(t *testing.T) {
	testTracks := []string{
		"{ \"latitude\" : 4.5343, \"longitude\" : 3.34324, \"driver_id\": 1 }",
		"{ \"latitude\" : 4.6424, \"longitude\" : 3.35232, \"driver_id\": 1 }",
		"{ \"latitude\" : 4.7534, \"longitude\" : 3.36534, \"driver_id\": 2 }",
		"{ \"latitude\" : 4.8442, \"longitude\" : 3.37321, \"driver_id\": 2 }",
		"{ \"latitude\" : 4.9224, \"longitude\" : 3.38234, \"driver_id\": 2 }",
	}

	server := New(NewRouter(), NewJSONAdapter(), "", 80)
	server.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTracRepoMemory()), "POST")

	testServer := httptest.NewServer(server.httpServer.Handler)
	defer testServer.Close()

	for _, track := range testTracks {
		resp, err := http.Post(testServer.URL+"/track/store", "application/json", bytes.NewBufferString(track))
		if err != nil {
			t.Errorf("Error requesting url <%s>", err)
		}

		// Check status code
		if got, expected := resp.StatusCode, http.StatusOK; got != expected {
			t.Errorf("Bad Status Code. Got <%v>, expecting <%v>", got, expected)
		}

		// Check content type
		if got, expected := resp.Header.Get("Content-Type"), "application/json"; got != expected {
			t.Errorf("Bad Content Type. Got <%s>, expecting <%s>", got, expected)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error reading body <%s>", err)
		}

		// Check that response is a valid json
		got := ""
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("Error decoding json response <%s>", err)
		}

		// Check response
		if expected := "OK"; got != expected {
			t.Errorf("Wrong body response for /track/store. Got <%s>, expecting <%s>", got, expected)
		}
		resp.Body.Close()
	}
}

func TestStoreAtrackWrongBody(t *testing.T) {
	testTracks := []string{
		"",
		"El Perro de san roque no tiene rabo",
		"{ \"latitude\" : , \"longitude\" : 3.36534, \"driver_id\": 2 }",
		"{ \"latitude\" : lala, \"longitude\" : 3.37321, \"driver_id\": 2 }",
	}

	server := New(NewRouter(), NewJSONAdapter(), "", 80)
	server.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTracRepoMemory()), "POST")

	testServer := httptest.NewServer(server.httpServer.Handler)
	defer testServer.Close()

	for _, track := range testTracks {
		resp, err := http.Post(testServer.URL+"/track/store", "application/json", bytes.NewBufferString(track))
		if err != nil {
			t.Errorf("Error requesting url <%s>", err)
		}

		// Check status code
		if got, expected := resp.StatusCode, http.StatusBadRequest; got != expected {
			t.Errorf("Bad Status Code. Got <%v>, expecting <%v>", got, expected)
		}

		// Check content type
		if got, expected := resp.Header.Get("Content-Type"), "application/json"; got != expected {
			t.Errorf("Bad Content Type. Got <%s>, expecting <%s>", got, expected)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error reading body <%s>", err)
		}

		// Check that response is a valid json
		got := ""
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("Error decoding json response <%s>", err)
		}

		resp.Body.Close()
	}
}
