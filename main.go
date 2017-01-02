package main

import (
	"math/rand"
	"net/http"
	"time"

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
	mainLogger := logger.WithField("func", "main")

	// Seed the random number generator, transaction IDs are random
	rand.Seed(time.Now().UnixNano())

	var spec specification
	err = envconfig.Process("APP", &spec)
	if err != nil {
		mainLogger.Fatalln(err)
	}
	mainLogger.Info(spec)

	s, err := NewServer()
	if err != nil {
		mainLogger.Fatalln(err)
	}
	s.log = logger

	r := mux.NewRouter()
	r.Handle("/stream/{dest}", s.GetStreamTracerouteHandler()).Methods("GET")
	r.Handle("/traceroute/{dest}", s.GetTracerouteHandler()).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("build")))

	mainLogger.Info("Starting API server")
	err = http.ListenAndServe(spec.Bind, r)
	if err != nil {
		mainLogger.Errorln(err)
	}
}
