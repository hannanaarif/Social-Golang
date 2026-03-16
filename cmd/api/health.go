package main

import (
	"net/http"
)

// healthCheckHandler godoc
// @Summary		Health Check
// @Description	Checks if the API is running
// @Tags			ops
// @Produce		json
// @Success		200	{object}	map[string]string
// @Failure		500	{object}	error
// @Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"message": "success",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		// log.Println(err)
		app.InternalServerError(w, r, err)
		// app.errorJSON(w,http.StatusInternalServerError,"err.Error")
	}
}
