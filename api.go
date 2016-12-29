package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Traceroute contains details of a traceroute
type Traceroute struct {
	Destination string `json:"destination"`
	Hops        []Hop  `json:"hops"`
}

// GetTracerouteHandler performs a live traceroute
func (s *Server) GetTracerouteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		s.log.WithField("func", "GetTraceroute")
		s.log.Infoln(r.Method, r.URL.Path, r.RemoteAddr)

		vars := mux.Vars(r)
		dest := vars["dest"]

		response := Traceroute{
			Destination: dest,
		}
		response.Hops, err = traceroute(dest)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
		}

		httpResponse(w, &response, http.StatusOK)
	})
}
