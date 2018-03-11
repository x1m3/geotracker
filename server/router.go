package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"io"
)

type Router struct {
	mux.Router
}

func NewRouter() *Router{
	router := &Router{}
	router.NotFoundHandler = router.RouteNotFoundHandler()
	return router
}

func (r *Router) RouteNotFoundHandler() http.HandlerFunc{
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusNotFound)
		resp.Header().Set("Content-Type", "text/html")
		io.WriteString(resp, "Not Found.")
	}
}


