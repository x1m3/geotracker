package Server

import (
	"net/http"
	"io"
	"encoding/json"
	"time"
	"log"
)

// Maximum time to read the full http request
const SERVER_HTTP_READTIMEOUT = 30 * time.Second

// Maximum time to write the full http request
const SERVER_HTTP_WRITETIMEOUT = 30 * time.Second

// Keep alive timeout. Time to close an idle connection if keep alive is enable
const SERVER_HTTP_IDLETIMEOUT = 5 * time.Second

// Time to wait for pending connections to finish when doing a shutdown
const TIMEOUT_SHUTDOWN_WAIT_PENDING_CONNS = 30 * time.Second


type Server struct {
	httpServer  *http.Server
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

	router.HandleFunc("/ping", server.ping).Methods("GET")
	return server
}

func (s *Server) Run() {
	log.Print("Starting server")
	err := s.httpServer.ListenAndServe()
	if err!=nil {
		log.Fatalf("Cannot start server. Reason <%s>", err)
	}
}




func (s *Server) encode(resp io.Writer, item interface{}) error {
	jsonEncoder := json.NewEncoder(resp)
	return jsonEncoder.Encode(item)
}

func (s *Server) ping(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	err := s.encode(resp, "pong")
	if err!=nil {
		resp.WriteHeader(http.StatusOK)
	}
}
