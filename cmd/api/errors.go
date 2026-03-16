package main

import (
	"net/http"
)

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err);
	app.logger.Error("Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Bad request: %s path: %s error: %s", r.Method, r.URL.Path, err);
	app.logger.Error("Bad request: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Not found: %s path: %s error: %s", r.Method, r.URL.Path, err);
	app.logger.Error("Not found: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, message)
}
