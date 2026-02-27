package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"message": "success",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		// log.Println(err)
		app.	InternalServerError(w, r, err)
		// writeJSONError(w,http.StatusInternalServerError,"err.Error")
	}
}


