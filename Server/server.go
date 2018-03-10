package Server

import (
	"github.com/gorilla/mux"
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

func New() *Server {
	server := &Server{}
	server.httpServer = &http.Server{
		Handler:      server.initRouter(),
		ReadTimeout:  SERVER_HTTP_READTIMEOUT,
		WriteTimeout: SERVER_HTTP_WRITETIMEOUT,
		IdleTimeout:  SERVER_HTTP_IDLETIMEOUT,
	}
	server.httpServer.SetKeepAlivesEnabled(true)
	return server
}

func (s *Server) Run() {
	log.Print("Starting server")
	err := s.httpServer.ListenAndServe()
	if err!=nil {
		log.Fatalf("Cannot start server. Reason <%s>", err)
	}
}

func (s *Server) initRouter() *mux.Router{
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusNotFound)
			resp.Header().Set("Content-Type", "text/html")
			io.WriteString(resp, "Not Found.")
		})

	router.HandleFunc("/ping", s.ping).Methods("GET")
	return router
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
