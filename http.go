package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorResponse is used when a json object needs to be returned with just an error
type errorResponse struct {
	Error string `json:"error"`
}

// httpResponse is a nice wrapper for sending JSON responses
func httpResponse(w http.ResponseWriter, v interface{}, status int) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, string(body))
	return nil
}

// streamResponse is a nice wrapper for sending multiple JSON responses
func streamResponse(w http.ResponseWriter, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(body))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}
