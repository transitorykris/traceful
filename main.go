package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

// this service's configuration
type specification struct {
	Bind string `envconfig:"bind" default:":8080"`
}

func main() {
	var err error

	// Set up our logging options
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Level = logrus.DebugLevel

	var spec specification
	err = envconfig.Process("APP", &spec)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Info(spec)

	s, err := NewServer()
	if err != nil {
		logger.Fatalln(err)
	}
	s.log = logger

	r := mux.NewRouter()
	r.Handle("/traceroute/{dest}", s.GetTracerouteHandler()).Methods("GET")

	logger.Info("Starting API server")
	err = http.ListenAndServe(spec.Bind, r)
	if err != nil {
		logger.Errorln(err)
	}
}
