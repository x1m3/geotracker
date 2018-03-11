package server

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func TestServer_Run(t *testing.T) {
	server := New(NewRouter(), NewJSONAdapter())

	testServer := httptest.NewServer(server.httpServer.Handler)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/ping")
	if err != nil {
		t.Errorf("Error pinging server. <%s>", err)
	}

	if got, expected := resp.StatusCode, http.StatusOK; got != expected {
		t.Errorf("Bad Status Code. Got <%v>, expecting <%v>", got, expected)
	}
	if got, expected := resp.Header.Get("Content-Type"), "application/json"; got!=expected {
		t.Errorf("Bad Content Type. Got <%s>, expecting <%s>", got, expected)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		t.Errorf("Error reading body <%s>", err)
	}

	got:=""
	if err := json.Unmarshal(body,&got); err!=nil {
		t.Errorf("Error decoding json response <%s>", err)
	}

	if expected := "pong"; got!=expected {
		t.Errorf("Wrong body response for /ping. Got <%s>, expecting <%s>", got, expected)
	}
}
