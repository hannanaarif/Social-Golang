package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 //1mb
	//http.MaxBytesReader wraps the request body. If a client sends more than 1MB,
	// the reader will stop and return an error. This prevents "Denial of Service" attacks where
	// someone sends a massive file to crash the server.
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	/*
		ensures that if the client sends fields that aren't defined in your Go struct (e.g., RegisterUserPayload), the decoding fails. This is great for catching typos or preventing users from sending unexpected data.
	*/
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(w, status, &envelope{Error: message})
}


