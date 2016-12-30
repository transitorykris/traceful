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

// GetTracerouteHandler performs a live traceroute
func (s *Server) GetTracerouteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.WithField("func", "GetTraceroute")
		s.log.Infoln(r.RemoteAddr, r.Method, r.URL.Path, r.URL.Query())

		vars := mux.Vars(r)
		dest := vars["dest"]

		var opts []TraceOpt

		if r.URL.Query().Get("hops") != "" {
			hops, err := strconv.Atoi(r.URL.Query().Get("hops"))
			if err != nil {
				httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
			opts = append(opts, HopsOpt(hops))
		}

		if r.URL.Query().Get("retries") != "" {
			retries, err := strconv.Atoi(r.URL.Query().Get("retries"))
			if err != nil {
				httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
			opts = append(opts, RetriesOpt(retries))
		}

		if r.URL.Query().Get("timeout") != "" {
			timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))
			if err != nil {
				httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
			opts = append(opts, TimeoutOpt(timeout))
		}

		if r.URL.Query().Get("port") != "" {
			port, err := strconv.Atoi(r.URL.Query().Get("port,"))
			if err != nil {
				httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
			opts = append(opts, PortOpt(port))
		}

		if r.URL.Query().Get("size") != "" {
			size, err := strconv.Atoi(r.URL.Query().Get("size"))
			if err != nil {
				httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
			opts = append(opts, SizeOpt(size))
		}

		response := Traceroute{
			Destination: dest,
		}
		var err error
		response.Hops, err = traceroute(dest, opts...)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		httpResponse(w, &response, http.StatusOK)
	})
}
