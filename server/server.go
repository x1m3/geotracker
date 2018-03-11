package server

import (
	"net/http"
	"io"
	"encoding/json"
	"time"
	"log"
	"bytes"
	"fmt"
	"github.com/x1m3/geotracker/command"
)

// Maximum time to read the full http request
const SERVER_HTTP_READTIMEOUT = 30 * time.Second

// Maximum time to write the full http request
const SERVER_HTTP_WRITETIMEOUT = 30 * time.Second

// Keep alive timeout. Time to close an idle connection if keep alive is enable
const SERVER_HTTP_IDLETIMEOUT = 5 * time.Second

// Time to wait for pending connections to finish when doing a shutdown
const TIMEOUT_SHUTDOWN_WAIT_PENDING_CONNS = 30 * time.Second

type encoder interface {
	encode(resp io.Writer, item interface{}) error
}

type Server struct {
	httpServer *http.Server
}

func New(router *Router) *Server {
	server := &Server{}
	server.httpServer = &http.Server{
		Handler:      router,
		ReadTimeout:  SERVER_HTTP_READTIMEOUT,
		WriteTimeout: SERVER_HTTP_WRITETIMEOUT,
		IdleTimeout:  SERVER_HTTP_IDLETIMEOUT,
	}
	server.httpServer.SetKeepAlivesEnabled(true)

	router.HandleFunc("/ping", server.handle(server.ping)).Methods("GET")
	return server
}

func (s *Server) Run() {
	log.Print("Starting server")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot start server. Reason <%s>", err)
	}
}

func (s *Server) encode(resp io.Writer, item interface{}) error {
	jsonEncoder := json.NewEncoder(resp)
	return jsonEncoder.Encode(item)
}

// server.handle returns a function that satisfies http.HandlerFunc. It's purpose is to execute a command.Command and
// adapt command.Request and command.Response to an http request. This way, commands could be usable in other protocols, simple
// implementing a new server that would translate a command response to this other protocol
//
// 1) Decodes the body
// 2) Runs a command
// 3) Based on the response of the command it returns an error or it decodes the response and
//    writes the response in http format.
func (s *Server) handle(command command.Command) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		request:= make(map[string]interface{},0)

		// Write a proper content-type header
		resp.Header().Set("Content-Type", "application/json")

		// Decode the body.
		decoder := json.NewDecoder(req.Body)
		if decoder.More() {
			if err := decoder.Decode(&request); err != nil {
				if err != nil {
					msg := fmt.Sprintf("Error decoding body. <%s>", err)
					log.Printf(msg)
					http.Error(resp, msg, http.StatusBadRequest)
					return
				}
			}
		}

		// Run the command
		response, err := command(request)
		if err != nil {
			msg := fmt.Sprintf("Error executing command <%s>", err)
			log.Printf(msg)
			http.Error(resp, msg, http.StatusInternalServerError)
			return
		}

		// Decode the command response
		buff := bytes.Buffer{}
		err = s.encode(&buff, response)
		if err != nil {
			msg := fmt.Sprintf("Error encoding response <%s>", err)
			log.Printf(msg)
			http.Error(resp, msg, http.StatusInternalServerError)
			return
		}

		// Write the response to the caller
		resp.WriteHeader(http.StatusOK)
		buff.WriteTo(resp)
	}
}

func (s *Server) ping(req command.Request) (command.Response, error) {
	return "pong", nil
}
