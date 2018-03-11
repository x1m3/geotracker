package HTTPServer

import (
	"net/http"
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

type Server struct {
	httpServer      *http.Server
	protocolAdapter ProtocolAdapter
}

func New(router *Router, adapter ProtocolAdapter) *Server {
	server := &Server{protocolAdapter: adapter}
	server.httpServer = &http.Server{
		Handler:      router,
		ReadTimeout:  SERVER_HTTP_READTIMEOUT,
		WriteTimeout: SERVER_HTTP_WRITETIMEOUT,
		IdleTimeout:  SERVER_HTTP_IDLETIMEOUT,
	}
	server.httpServer.SetKeepAlivesEnabled(true)
	return server
}

func (s *Server) RegisterEndpoint(path string, command command.Command, method string) {
	s.httpServer.Handler.(*Router).HandleFunc(path, s.handle(command)).Methods(method)
}

func (s *Server) Run() {
	log.Print("Starting HTTPServer")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot start HTTPServer. Reason <%s>", err)
	}
}

// HTTPServer.handle returns a function that satisfies http.HandlerFunc. It's purpose is to execute a command.Command and
// adapt command.Request and command.Response to an http request. This way, commands could be usable in other protocols, simple
// implementing a new HTTPServer that would translate a command response to this other protocol
//
// 1) Decodes the body
// 2) Runs a command
// 3) Based on the response of the command it returns an error or it decodes the response and
//    writes the response in http format.
func (s *Server) handle(command command.Command) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		request := make(map[string]interface{}, 0)

		// Write a proper content-type header
		resp.Header().Set("Content-Type", s.protocolAdapter.ContentType())

		// Decode the body.
		request, err := s.protocolAdapter.Decode(req.Body)
		if err != nil {
			if err != nil {
				resp.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("Error decoding body. <%s>", err)
				log.Printf(msg)
				s.protocolAdapter.Encode(resp, msg)
				return
			}
		}

		// Run the command
		response, err := command.Call(request)
		if err != nil {
			resp.WriteHeader(s.mapError(err))
			msg := fmt.Sprintf("Error executing command <%s>", err)
			log.Printf(msg)
			s.protocolAdapter.Encode(resp, msg)
			return
		}

		// Decode the command response
		buff := bytes.Buffer{}
		err = s.protocolAdapter.Encode(&buff, response)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Error encoding response <%s>", err)
			log.Printf(msg)
			s.protocolAdapter.Encode(resp, msg)
			return
		}

		// Write the response to the caller
		resp.WriteHeader(http.StatusOK)
		buff.WriteTo(resp)
	}
}

func (s *Server) mapError(err error) int {
	switch err {
	case command.ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
