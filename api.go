package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// Traceroute contains details of a traceroute
type Traceroute struct {
	Destination string `json:"destination"`
	Hops        []Hop  `json:"hops"`
}

// Hop contains details of a hop in a traceroute
type Hop struct {
	TraceHop
	GeoIP
}

// GeoIP contains geoip details for a hop
type GeoIP struct {
	Country string `json:"country,omitempty"`
	ASN     int    `json:"asn,omitempty"`
}

// getGeoIP gets geoip data for an IP
func (s *Server) getGeoIP(ip string) (GeoIP, error) {
	geoip := GeoIP{}
	if s.geoIPURL == "" {
		return geoip, nil
	}
	resp, err := http.Get(fmt.Sprintf("%s/ip/%s", s.geoIPURL, ip))
	if err != nil {
		return geoip, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return geoip, err
	}
	err = json.Unmarshal(body, &geoip)
	return geoip, err
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
		log := s.log.WithFields(logrus.Fields{"func": "GetTraceroute", "id": rand.Int63()})
		log.WithFields(logrus.Fields{"remote": r.RemoteAddr, "method": r.Method, "path": r.URL.Path, "query": r.URL.Query()}).Infoln()

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
		traceHops, err := traceroute(dest, opts...)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
		var hop Hop
		for _, h := range traceHops {
			hop.TraceHop = h
			hop.GeoIP, _ = s.getGeoIP(h.Address)
			response.Hops = append(response.Hops, hop)
		}

		httpResponse(w, &response, http.StatusOK)
	})
}

// GetStreamTracerouteHandler performs a live traceroute
func (s *Server) GetStreamTracerouteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := s.log.WithFields(logrus.Fields{"func": "GetStreamTracerouteHandler", "id": rand.Int63()})
		log.WithFields(logrus.Fields{"remote": r.RemoteAddr, "method": r.Method, "path": r.URL.Path, "query": r.URL.Query()}).Infoln()

		vars := mux.Vars(r)
		dest := vars["dest"]

		opts, err := ParamsToOpts(r)
		if err != nil {
			httpResponse(w, &errorResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}

		ch := make(chan TraceHop, 0)
		done := make(chan bool)
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			httpResponse(w, errorResponse{Error: "Cannot stream"}, http.StatusInternalServerError)
			return
		}
		closeNotify := cn.CloseNotify()
		go func() {
			w.Header().Set("Content-Type", "application/stream+json")
			for {
				select {
				case hop, ok := <-ch:
					if !ok {
						log.Errorln("problem completing traceroute")
						streamResponse(w, &errorResponse{Error: "problem completing traceroute"})
						return
					}
					geoip, _ := s.getGeoIP(hop.Address)
					resp := Hop{
						TraceHop: hop,
						GeoIP:    geoip,
					}
					streamResponse(w, resp)
				case <-done:
					return
				case <-closeNotify:
					return
				}
			}
		}()

		err = liveTraceroute(dest, ch, done, opts...)
		if err != nil {
			log.Errorln(err)
			streamResponse(w, &errorResponse{Error: err.Error()})
			return
		}
		return
	})
}
