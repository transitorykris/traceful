package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Traceroute contains details of a traceroute
type Traceroute struct {
	Destination string `json:"destination"`
}

// GetTracerouteHandler performs a live traceroute
func (s *Server) GetTracerouteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.WithField("func", "GetTraceroute")
		s.log.Infoln(r.Method, r.URL.Path, r.RemoteAddr)

		vars := mux.Vars(r)
		dest := vars["dest"]

		httpResponse(w, &Traceroute{Destination: dest}, http.StatusOK)
	})
}
