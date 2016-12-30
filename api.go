package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Traceroute contains details of a traceroute
type Traceroute struct {
	Destination string `json:"destination"`
	Hops        []Hop  `json:"hops"`
}

// ParamsToOpts converts the query params into TraceOpts
func ParamsToOpts(r *http.Request) ([]TraceOpt, error) {
	var opts []TraceOpt

	if r.URL.Query().Get("hops") != "" {
		hops, err := strconv.Atoi(r.URL.Query().Get("hops"))
		if err != nil {
			return opts, err
		}
		opts = append(opts, HopsOpt(hops))
	}

	if r.URL.Query().Get("retries") != "" {
		retries, err := strconv.Atoi(r.URL.Query().Get("retries"))
		if err != nil {
			return opts, err
		}
		opts = append(opts, RetriesOpt(retries))
	}

	if r.URL.Query().Get("timeout") != "" {
		timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))
		if err != nil {
			return opts, err
		}
		opts = append(opts, TimeoutOpt(timeout))
	}

	if r.URL.Query().Get("size") != "" {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			return opts, err
		}
		opts = append(opts, SizeOpt(size))
	}

	return opts, nil
}

// GetTracerouteHandler performs a live traceroute
func (s *Server) GetTracerouteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.WithField("func", "GetTraceroute")
		s.log.Infoln(r.RemoteAddr, r.Method, r.URL.Path, r.URL.Query())

		vars := mux.Vars(r)
		dest := vars["dest"]

		opts, err := ParamsToOpts(r)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		response := Traceroute{
			Destination: dest,
		}
		response.Hops, err = traceroute(dest, opts...)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		httpResponse(w, &response, http.StatusOK)
	})
}
